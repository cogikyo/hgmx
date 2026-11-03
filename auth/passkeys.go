package auth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncbor"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/pg"
)

const (
	login    = "login"
	register = "register"
)

type passkey struct {
	ID         []byte
	UserID     pg.ID
	Name       string
	Credential pg.JSON[webauthn.Credential]
	Created    time.Time
	Used       *time.Time
}

func (passkey) TableName() string { return "auth_passkeys" }

type ceremony struct {
	ID        []byte
	Kind      string
	Session   pg.JSON[webauthn.SessionData]
	UserID    *pg.ID
	Library   string
	Expires   time.Time
	ClientID  *string
	Origin    *string
	SessionID *pg.ID
}

func (ceremony) TableName() string { return "auth_ceremonies" }

var (
	errUnknown    = errors.New("auth: passkey not registered")
	errOrigin     = errors.New("auth: passkey response from another origin")
	errUnverified = errors.New("auth: passkey could not be verified")
)

func (a *Auth) verifier(c caller) *webauthn.WebAuthn {
	if c.app == nil {
		return a.wa
	}
	return c.app.wa
}

func (c caller) proves(origin string) bool {
	return c.app == nil || slices.Contains(c.app.origins, origin)
}

func (a *Auth) saveCeremony(ctx context.Context, c caller, kind string, user, from *pg.ID, s *webauthn.SessionData) (string, error) {
	db := a.db.WithContext(ctx)
	if err := sweep(db, a.library); err != nil {
		return "", err
	}
	t := token()
	row := ceremony{ID: hash(t), Kind: kind, Session: pg.JSON[webauthn.SessionData]{V: *s}, UserID: user, Library: a.library, Expires: time.Now().Add(ceremonyLife)}
	if c.app != nil {
		row.ClientID, row.Origin, row.SessionID = c.client(), &c.host.origin, from
	}
	return t, db.Create(&row).Error
}

func takeCeremony(db *gorm.DB, c caller, t, kind, library string) (ceremony, bool, error) {
	var row ceremony
	q, args := "DELETE FROM auth_ceremonies WHERE id = ? AND session->>'origin' = ? AND client_id IS NULL RETURNING *", []any{hash(t), c.host.origin}
	if c.app != nil {
		q, args = "DELETE FROM auth_ceremonies WHERE id = ? AND origin = ? AND client_id = ? RETURNING *", append(args, c.app.id)
	}
	if err := db.Raw(q, args...).Scan(&row).Error; err != nil {
		return row, false, err
	}
	return row, row.ID != nil && row.Kind == kind && row.Library == library && time.Now().Before(row.Expires), nil
}

func (a *Auth) ceremony(w http.ResponseWriter, r *http.Request, h *host, kind string) (ceremony, bool) {
	cookie, err := r.Cookie(ceremonyCookie)
	if err != nil {
		return ceremony{}, false
	}
	setCookie(w, ceremonyCookie, "", -time.Second)
	c, ok, err := takeCeremony(a.db.WithContext(r.Context()), browser(r, h), cookie.Value, kind, a.library)
	if err != nil {
		panic(err)
	}
	return c, ok
}

func (a *Auth) beginLogin(ctx context.Context, c caller) (*protocol.CredentialAssertion, string, error) {
	var opts []webauthn.LoginOption
	if c.app == nil {
		opts = append(opts, webauthn.WithLoginOrigin(c.host.origin))
	}
	options, s, err := a.verifier(c).BeginDiscoverableLogin(opts...)
	if err != nil {
		return nil, "", err
	}
	t, err := a.saveCeremony(ctx, c, login, nil, nil, s)
	return options, t, err
}

func (a *Auth) loginBegin(w http.ResponseWriter, r *http.Request, h *host) {
	c := browser(r, h)
	if !h.begins.allow(c.ip) {
		problem(w, http.StatusTooManyRequests, "Too many sign-in attempts. Try again in a few minutes.")
		return
	}
	options, t, err := a.beginLogin(r.Context(), c)
	if err != nil {
		panic(err)
	}
	setCookie(w, ceremonyCookie, t, ceremonyLife)
	writeJSON(w, options)
}

func (a *Auth) passkeyLogin(ctx context.Context, c caller, cer ceremony, parsed *protocol.ParsedCredentialAssertionData) (string, error) {
	if !c.proves(parsed.Response.CollectedClientData.Origin) {
		return "", errOrigin
	}
	var t string
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var u *holder
		var load error
		_, cred, err := a.verifier(c).ValidatePasskeyLogin(func(id, handle []byte) (webauthn.User, error) {
			// Lock user before credential to avoid a cycle with account deletion.
			// KEY SHARE permits passkey removal's NO KEY UPDATE user lock.
			// Fill the username before KEY SHARE to avoid concurrent lock upgrades.
			if load = claim(tx, "handle", handle); load != nil {
				return nil, load
			}
			if load = tx.Exec("SELECT 1 FROM auth_users WHERE handle = ? FOR KEY SHARE", handle).Error; load != nil {
				return nil, load
			}
			if load = tx.Exec("SELECT 1 FROM auth_passkeys WHERE id = ? FOR NO KEY UPDATE", id).Error; load != nil {
				return nil, load
			}
			u, load = loadHolder(tx, "handle", handle)
			return u, load
		}, cer.Session.V, parsed)
		if load != nil && !errors.Is(load, pg.ErrNotFound) {
			return load
		}
		if err != nil {
			return errUnknown
		}
		// Validation and row-lock waits can outlive the initial expiry check.
		if !time.Now().Before(cer.Expires) {
			return errStale
		}
		if !c.host.admits(u.User) {
			return errDenied
		}
		t, err = finishLogin(tx, c, u.ID, cred)
		if errors.Is(err, pg.ErrNotFound) {
			return errUnknown
		}
		return err
	})
	return t, err
}

func (a *Auth) loginFinish(w http.ResponseWriter, r *http.Request, h *host) {
	c, ok := a.ceremony(w, r, h, login)
	if !ok {
		problem(w, http.StatusBadRequest, "Sign-in took too long. Try again.")
		return
	}
	parsed, err := protocol.ParseCredentialRequestResponseBody(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil {
		problem(w, http.StatusBadRequest, "That passkey response was not readable.")
		return
	}
	t, err := a.passkeyLogin(r.Context(), browser(r, h), c, parsed)
	if errors.Is(err, errUnknown) {
		problem(w, http.StatusUnauthorized, "That passkey is not registered here.")
		return
	}
	if errors.Is(err, errStale) {
		problem(w, http.StatusBadRequest, "Sign-in took too long. Try again.")
		return
	}
	if errors.Is(err, errDenied) {
		problem(w, http.StatusForbidden, "This account can't sign in here.")
		return
	}
	if err != nil {
		panic(err)
	}
	setCookie(w, sessionCookie, t, sessionLife)
	w.WriteHeader(http.StatusNoContent)
}

func finishLogin(tx *gorm.DB, c caller, user pg.ID, cred *webauthn.Credential) (string, error) {
	res := tx.Model(&passkey{}).Where("id = ? AND user_id = ?", cred.ID, user).
		Updates(map[string]any{"credential": pg.JSON[webauthn.Credential]{V: *cred}, "used": time.Now()})
	if res.Error != nil {
		return "", res.Error
	}
	if res.RowsAffected != 1 {
		return "", pg.ErrNotFound
	}
	return startSession(tx, c, user)
}

func (a *Auth) recentSession(w http.ResponseWriter, r *http.Request, h *host) (session, bool) {
	s, ok := a.session(r, h)
	if !ok {
		problem(w, http.StatusUnauthorized, "Your session ended. Sign in again.")
		return s, false
	}
	if !h.admits(s.user) {
		problem(w, http.StatusForbidden, "This account can't change passkeys here.")
		return s, false
	}
	if !s.recent() {
		problem(w, http.StatusForbidden, "Sign in again to change passkeys.")
		return s, false
	}
	return s, true
}

func (a *Auth) beginAdd(ctx context.Context, c caller, s session) (*protocol.CredentialCreation, string, error) {
	u, err := loadHolder(a.db.WithContext(ctx), "id", s.user.ID)
	if err != nil {
		return nil, "", err
	}
	var exclude []protocol.CredentialDescriptor
	for _, k := range u.creds {
		if compatible(k, a.rp) {
			exclude = append(exclude, k.Descriptor())
		}
	}
	opts := []webauthn.RegistrationOption{webauthn.WithExclusions(exclude)}
	if c.app == nil {
		opts = append(opts, webauthn.WithRegistrationOrigin(c.host.origin))
	}
	options, sd, err := a.verifier(c).BeginRegistration(u, opts...)
	if err != nil {
		return nil, "", err
	}
	t, err := a.saveCeremony(ctx, c, register, &u.ID, &s.id, sd)
	return options, t, err
}

func (a *Auth) addBegin(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.recentSession(w, r, h)
	if !ok {
		return
	}
	c := browser(r, h)
	if !h.begins.allow(c.ip) {
		problem(w, http.StatusTooManyRequests, "Too many attempts. Try again in a few minutes.")
		return
	}
	options, t, err := a.beginAdd(r.Context(), c, s)
	if err != nil {
		panic(err)
	}
	setCookie(w, ceremonyCookie, t, ceremonyLife)
	writeJSON(w, options)
}

func (a *Auth) finishAdd(ctx context.Context, c caller, s session, cer ceremony, parsed *protocol.ParsedCredentialCreationData) (string, User, error) {
	if !c.proves(parsed.Response.CollectedClientData.Origin) {
		return "", User{}, errOrigin
	}
	db := a.db.WithContext(ctx)
	holder, err := loadHolder(db, "id", s.user.ID)
	if err != nil {
		return "", User{}, err
	}
	cred, err := a.verifier(c).CreateCredential(holder, cer.Session.V, parsed)
	if err != nil {
		return "", User{}, errUnverified
	}
	var name string
	var u User
	err = db.Transaction(func(tx *gorm.DB) error {
		var err error
		if u, err = hold(tx, s.user.ID); err != nil {
			return err
		}
		// Reset may have revoked the session while the ceremony was in flight.
		cur, err := live(tx, c, s.user.ID)
		if err != nil {
			return err
		}
		if !cur.recent() || !time.Now().Before(cer.Expires) || cer.SessionID != nil && cur.id != *cer.SessionID {
			return errStale
		}
		var n int64
		if err := tx.Model(&passkey{}).Where("user_id = ?", s.user.ID).Count(&n).Error; err != nil {
			return err
		}
		name = fmt.Sprintf("Passkey %d", n+1)
		return tx.Create(&passkey{ID: cred.ID, UserID: s.user.ID, Name: name, Credential: pg.JSON[webauthn.Credential]{V: *cred}, Created: time.Now()}).Error
	})
	return name, u, err
}

func (a *Auth) addFinish(w http.ResponseWriter, r *http.Request, h *host) {
	c, fresh := a.ceremony(w, r, h, register)
	parsed, parseErr := protocol.ParseCredentialCreationResponseBody(http.MaxBytesReader(w, r.Body, 64<<10))
	s, ok := a.recentSession(w, r, h)
	if !ok {
		return
	}
	if !fresh || c.UserID == nil || *c.UserID != s.user.ID {
		problem(w, http.StatusBadRequest, "Adding the passkey took too long. Try again.")
		return
	}
	if parseErr != nil {
		problem(w, http.StatusBadRequest, "That passkey response was not readable.")
		return
	}
	name, u, err := a.finishAdd(r.Context(), browser(r, h), s, c, parsed)
	if errors.Is(err, errUnverified) {
		problem(w, http.StatusBadRequest, "That passkey could not be verified.")
		return
	}
	if errors.Is(err, errDenied) {
		problem(w, http.StatusForbidden, "This account can't change passkeys here.")
		return
	}
	if errors.Is(err, errStale) || errors.Is(err, pg.ErrNotFound) {
		problem(w, http.StatusForbidden, "Your sign-in changed while adding the passkey. Sign in again.")
		return
	}
	if errors.Is(err, pg.ErrDuplicate) {
		problem(w, http.StatusConflict, "That passkey is already saved.")
		return
	}
	if err != nil {
		panic(err)
	}
	a.notify(a.passkeyMail(r.Context(), h, u.Email, name, true))
	w.WriteHeader(http.StatusNoContent)
}

func passkeyID(r *http.Request) ([]byte, bool) {
	return parsePasskeyID(r.PathValue("id"))
}

func parsePasskeyID(raw string) ([]byte, bool) {
	id, err := base64.RawURLEncoding.DecodeString(raw)
	return id, err == nil && len(id) > 0
}

func validPasskeyName(name string) bool { return name != "" && len([]rune(name)) <= 64 }

func (a *Auth) renamePasskey(ctx context.Context, c caller, user pg.ID, id []byte, name string) error {
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := hold(tx, user); err != nil {
			return err
		}
		if _, err := live(tx, c, user); err != nil {
			return err
		}
		res := tx.Model(&passkey{}).Where("id = ? AND user_id = ?", id, user).Update("name", name)
		if res.Error == nil && res.RowsAffected == 0 {
			return errMissing
		}
		return res.Error
	})
}

func (a *Auth) rename(w http.ResponseWriter, r *http.Request, h *host) {
	name := strings.TrimSpace(r.FormValue("name"))
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	id, ok := passkeyID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if !validPasskeyName(name) {
		a.accountPage(w, r, h, s, http.StatusUnprocessableEntity, "Passkey names need 1 to 64 characters.")
		return
	}
	err := a.renamePasskey(r.Context(), browser(r, h), s.user.ID, id, name)
	if errors.Is(err, errMissing) {
		http.NotFound(w, r)
		return
	}
	if a.refused(w, r, h, s, err) {
		return
	}
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func (a *Auth) removePasskey(ctx context.Context, c caller, user pg.ID, id []byte) (string, User, error) {
	var name string
	var u User
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		if u, err = hold(tx, user); err != nil {
			return err
		}
		if _, err := live(tx, c, user); err != nil {
			return err
		}
		var k passkey
		err = tx.Where("id = ? AND user_id = ?", id, user).Take(&k).Error
		if errors.Is(err, pg.ErrNotFound) {
			return errMissing
		}
		if err != nil {
			return err
		}
		n, err := a.methods(tx, c.host, u)
		if err != nil {
			return err
		}
		if compatible(k.Credential.V, a.rp) && n <= 1 {
			return errLast
		}
		name = k.Name
		return tx.Delete(&passkey{}, "id = ? AND user_id = ?", id, user).Error
	})
	return name, u, err
}

func (a *Auth) remove(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	if !s.recent() {
		a.accountPage(w, r, h, s, http.StatusForbidden, "Sign in again to remove a passkey.")
		return
	}
	id, ok := passkeyID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	name, u, err := a.removePasskey(r.Context(), browser(r, h), s.user.ID, id)
	if errors.Is(err, errMissing) {
		http.NotFound(w, r)
		return
	}
	if errors.Is(err, errLast) {
		a.accountPage(w, r, h, s, http.StatusConflict, "This is your only way to sign in. Add another passkey before removing it.")
		return
	}
	if a.refused(w, r, h, s, err) {
		return
	}
	a.notify(a.passkeyMail(r.Context(), h, u.Email, name, false))
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func encodeID(id []byte) string { return base64.RawURLEncoding.EncodeToString(id) }

func compatible(c webauthn.Credential, rp string) bool {
	var obj protocol.AttestationObject
	if err := webauthncbor.Unmarshal(c.Attestation.Object, &obj); err != nil || len(obj.RawAuthData) < 32 {
		return false
	}
	sum := sha256.Sum256([]byte(rp))
	return bytes.Equal(obj.RawAuthData[:32], sum[:])
}
