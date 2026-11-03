package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
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
	ID      []byte
	Kind    string
	Session pg.JSON[webauthn.SessionData]
	UserID  *pg.ID
	Library string
	Expires time.Time
}

func (ceremony) TableName() string { return "auth_ceremonies" }

func (a *Auth) saveCeremony(w http.ResponseWriter, r *http.Request, kind string, user *pg.ID, s *webauthn.SessionData) error {
	db := a.db.WithContext(r.Context())
	if err := sweep(db, a.library); err != nil {
		return err
	}
	t := token()
	err := db.Create(&ceremony{ID: hash(t), Kind: kind, Session: pg.JSON[webauthn.SessionData]{V: *s}, UserID: user, Library: a.library, Expires: time.Now().Add(ceremonyLife)}).Error
	if err != nil {
		return err
	}
	setCookie(w, ceremonyCookie, t, ceremonyLife)
	return nil
}

func takeCeremony(db *gorm.DB, h *host, t, kind, library string) (ceremony, bool, error) {
	var c ceremony
	err := db.Raw("DELETE FROM auth_ceremonies WHERE id = ? AND session->>'origin' = ? RETURNING *", hash(t), h.origin).Scan(&c).Error
	if err != nil {
		return c, false, err
	}
	return c, c.ID != nil && c.Kind == kind && c.Library == library && time.Now().Before(c.Expires), nil
}

func (a *Auth) ceremony(w http.ResponseWriter, r *http.Request, h *host, kind string) (ceremony, bool) {
	cookie, err := r.Cookie(ceremonyCookie)
	if err != nil {
		return ceremony{}, false
	}
	setCookie(w, ceremonyCookie, "", -time.Second)
	c, ok, err := takeCeremony(a.db.WithContext(r.Context()), h, cookie.Value, kind, a.library)
	if err != nil {
		panic(err)
	}
	return c, ok
}

func (a *Auth) loginBegin(w http.ResponseWriter, r *http.Request, h *host) {
	if !h.begins.allow(r) {
		problem(w, http.StatusTooManyRequests, "Too many sign-in attempts. Try again in a few minutes.")
		return
	}
	options, s, err := a.wa.BeginDiscoverableLogin(webauthn.WithLoginOrigin(h.origin))
	if err != nil {
		panic(err)
	}
	if err := a.saveCeremony(w, r, login, nil, s); err != nil {
		panic(err)
	}
	writeJSON(w, options)
}

var errUnknown = errors.New("auth: passkey not registered")

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
	var t string
	err = a.db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
		var u *holder
		var load error
		_, cred, err := a.wa.ValidatePasskeyLogin(func(id, handle []byte) (webauthn.User, error) {
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
		}, c.Session.V, parsed)
		if load != nil && !errors.Is(load, pg.ErrNotFound) {
			return load
		}
		if err != nil {
			return errUnknown
		}
		// Validation and row-lock waits can outlive the initial expiry check.
		if !time.Now().Before(c.Expires) {
			return errStale
		}
		if !h.admits(u.User) {
			return errDenied
		}
		t, err = finishLogin(tx, r, h, u.ID, cred)
		if errors.Is(err, pg.ErrNotFound) {
			return errUnknown
		}
		return err
	})
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

func finishLogin(tx *gorm.DB, r *http.Request, h *host, user pg.ID, cred *webauthn.Credential) (string, error) {
	res := tx.Model(&passkey{}).Where("id = ? AND user_id = ?", cred.ID, user).
		Updates(map[string]any{"credential": pg.JSON[webauthn.Credential]{V: *cred}, "used": time.Now()})
	if res.Error != nil {
		return "", res.Error
	}
	if res.RowsAffected != 1 {
		return "", pg.ErrNotFound
	}
	return startSession(tx, r, h, user)
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

func (a *Auth) addBegin(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.recentSession(w, r, h)
	if !ok {
		return
	}
	if !h.begins.allow(r) {
		problem(w, http.StatusTooManyRequests, "Too many attempts. Try again in a few minutes.")
		return
	}
	u, err := loadHolder(a.db.WithContext(r.Context()), "id", s.user.ID)
	if err != nil {
		panic(err)
	}
	var exclude []protocol.CredentialDescriptor
	for _, c := range u.creds {
		if compatible(c, a.rp) {
			exclude = append(exclude, c.Descriptor())
		}
	}
	options, sd, err := a.wa.BeginRegistration(u, webauthn.WithExclusions(exclude), webauthn.WithRegistrationOrigin(h.origin))
	if err != nil {
		panic(err)
	}
	if err := a.saveCeremony(w, r, register, &u.ID, sd); err != nil {
		panic(err)
	}
	writeJSON(w, options)
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
	db := a.db.WithContext(r.Context())
	holder, err := loadHolder(db, "id", s.user.ID)
	if err != nil {
		panic(err)
	}
	cred, err := a.wa.CreateCredential(holder, c.Session.V, parsed)
	if err != nil {
		problem(w, http.StatusBadRequest, "That passkey could not be verified.")
		return
	}
	var name string
	var u User
	err = db.Transaction(func(tx *gorm.DB) error {
		var err error
		if u, err = hold(tx, s.user.ID); err != nil {
			return err
		}
		// Reset may have revoked the session while the ceremony was in flight.
		cur, err := live(tx, r, h, s.user.ID)
		if err != nil {
			return err
		}
		if !cur.recent() || !time.Now().Before(c.Expires) {
			return errStale
		}
		var n int64
		if err := tx.Model(&passkey{}).Where("user_id = ?", s.user.ID).Count(&n).Error; err != nil {
			return err
		}
		name = fmt.Sprintf("Passkey %d", n+1)
		return tx.Create(&passkey{ID: cred.ID, UserID: s.user.ID, Name: name, Credential: pg.JSON[webauthn.Credential]{V: *cred}, Created: time.Now()}).Error
	})
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
	id, err := base64.RawURLEncoding.DecodeString(r.PathValue("id"))
	return id, err == nil && len(id) > 0
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
	if name == "" || len([]rune(name)) > 64 {
		a.accountPage(w, r, h, s, http.StatusUnprocessableEntity, "Passkey names need 1 to 64 characters.")
		return
	}
	err := a.db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
		if _, err := hold(tx, s.user.ID); err != nil {
			return err
		}
		if _, err := live(tx, r, h, s.user.ID); err != nil {
			return err
		}
		res := tx.Model(&passkey{}).Where("id = ? AND user_id = ?", id, s.user.ID).Update("name", name)
		if res.Error == nil && res.RowsAffected == 0 {
			return errMissing
		}
		return res.Error
	})
	if errors.Is(err, errMissing) {
		http.NotFound(w, r)
		return
	}
	if a.refused(w, r, h, s, err) {
		return
	}
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
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
	var name string
	var u User
	last := false
	err := a.db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
		var err error
		if u, err = hold(tx, s.user.ID); err != nil {
			return err
		}
		if _, err := live(tx, r, h, s.user.ID); err != nil {
			return err
		}
		var k passkey
		err = tx.Where("id = ? AND user_id = ?", id, s.user.ID).Take(&k).Error
		if errors.Is(err, pg.ErrNotFound) {
			return errMissing
		}
		if err != nil {
			return err
		}
		n, err := a.methods(tx, h, u)
		if err != nil {
			return err
		}
		if compatible(k.Credential.V, a.rp) && n <= 1 {
			last = true
			return nil
		}
		name = k.Name
		return tx.Delete(&passkey{}, "id = ? AND user_id = ?", id, s.user.ID).Error
	})
	if errors.Is(err, errMissing) {
		http.NotFound(w, r)
		return
	}
	if a.refused(w, r, h, s, err) {
		return
	}
	if last {
		a.accountPage(w, r, h, s, http.StatusConflict, "This is your only way to sign in. Add another passkey before removing it.")
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
