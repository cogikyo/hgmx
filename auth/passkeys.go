package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
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

func takeCeremony(db *gorm.DB, t, kind, library string) (ceremony, bool, error) {
	var c ceremony
	err := db.Raw("DELETE FROM auth_ceremonies WHERE id = ? RETURNING *", hash(t)).Scan(&c).Error
	if err != nil {
		return c, false, err
	}
	return c, c.ID != nil && c.Kind == kind && c.Library == library && time.Now().Before(c.Expires), nil
}

func (a *Auth) ceremony(w http.ResponseWriter, r *http.Request, kind string) (ceremony, bool) {
	cookie, err := r.Cookie(ceremonyCookie)
	if err != nil {
		return ceremony{}, false
	}
	setCookie(w, ceremonyCookie, "", -time.Second)
	c, ok, err := takeCeremony(a.db.WithContext(r.Context()), cookie.Value, kind, a.library)
	if err != nil {
		panic(err)
	}
	return c, ok
}

func (a *Auth) loginBegin(w http.ResponseWriter, r *http.Request) {
	if !a.begins.allow(r) {
		problem(w, http.StatusTooManyRequests, "Too many sign-in attempts. Try again in a few minutes.")
		return
	}
	options, s, err := a.wa.BeginDiscoverableLogin()
	if err != nil {
		panic(err)
	}
	if err := a.saveCeremony(w, r, login, nil, s); err != nil {
		panic(err)
	}
	writeJSON(w, options)
}

var errUnknown = errors.New("auth: passkey not registered")

func (a *Auth) loginFinish(w http.ResponseWriter, r *http.Request) {
	c, ok := a.ceremony(w, r, login)
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
		var h *holder
		var load error
		_, cred, err := a.wa.ValidatePasskeyLogin(func(id, handle []byte) (webauthn.User, error) {
			// Lock user before credential to avoid a cycle with account deletion.
			// KEY SHARE permits passkey removal's NO KEY UPDATE user lock.
			if load = tx.Exec("SELECT 1 FROM auth_users WHERE handle = ? FOR KEY SHARE", handle).Error; load != nil {
				return nil, load
			}
			if load = tx.Exec("SELECT 1 FROM auth_passkeys WHERE id = ? FOR NO KEY UPDATE", id).Error; load != nil {
				return nil, load
			}
			h, load = loadHolder(tx, "handle", handle)
			return h, load
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
		t, err = finishLogin(tx, r, h.ID, cred)
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
	if err != nil {
		panic(err)
	}
	setCookie(w, sessionCookie, t, sessionLife)
	w.WriteHeader(http.StatusNoContent)
}

func finishLogin(tx *gorm.DB, r *http.Request, user pg.ID, cred *webauthn.Credential) (string, error) {
	res := tx.Model(&passkey{}).Where("id = ? AND user_id = ?", cred.ID, user).
		Updates(map[string]any{"credential": pg.JSON[webauthn.Credential]{V: *cred}, "used": time.Now()})
	if res.Error != nil {
		return "", res.Error
	}
	if res.RowsAffected != 1 {
		return "", pg.ErrNotFound
	}
	return startSession(tx, r, user)
}

func (a *Auth) recentSession(w http.ResponseWriter, r *http.Request) (session, bool) {
	s, ok := a.session(r)
	if !ok {
		problem(w, http.StatusUnauthorized, "Your session ended. Sign in again.")
		return s, false
	}
	if !s.recent() {
		problem(w, http.StatusForbidden, "Sign in again to change passkeys.")
		return s, false
	}
	return s, true
}

func (a *Auth) addBegin(w http.ResponseWriter, r *http.Request) {
	s, ok := a.recentSession(w, r)
	if !ok {
		return
	}
	if !a.begins.allow(r) {
		problem(w, http.StatusTooManyRequests, "Too many attempts. Try again in a few minutes.")
		return
	}
	h, err := loadHolder(a.db.WithContext(r.Context()), "id", s.user.ID)
	if err != nil {
		panic(err)
	}
	exclude := make([]protocol.CredentialDescriptor, len(h.creds))
	for i, c := range h.creds {
		exclude[i] = c.Descriptor()
	}
	options, sd, err := a.wa.BeginRegistration(h, webauthn.WithExclusions(exclude))
	if err != nil {
		panic(err)
	}
	if err := a.saveCeremony(w, r, register, &h.ID, sd); err != nil {
		panic(err)
	}
	writeJSON(w, options)
}

func (a *Auth) addFinish(w http.ResponseWriter, r *http.Request) {
	c, fresh := a.ceremony(w, r, register)
	parsed, parseErr := protocol.ParseCredentialCreationResponseBody(http.MaxBytesReader(w, r.Body, 64<<10))
	s, ok := a.recentSession(w, r)
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
	h, err := loadHolder(db, "id", s.user.ID)
	if err != nil {
		panic(err)
	}
	cred, err := a.wa.CreateCredential(h, c.Session.V, parsed)
	if err != nil {
		problem(w, http.StatusBadRequest, "That passkey could not be verified.")
		return
	}
	var name string
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, h.Email); err != nil {
			return err
		}
		// Reset may have revoked the session while the ceremony was in flight.
		cur, ok, err := current(tx, r, h.ID)
		if err != nil {
			return err
		}
		if !ok || !cur.recent() || !time.Now().Before(c.Expires) {
			return errStale
		}
		var n int64
		if err := tx.Model(&passkey{}).Where("user_id = ?", h.ID).Count(&n).Error; err != nil {
			return err
		}
		name = fmt.Sprintf("Passkey %d", n+1)
		return tx.Create(&passkey{ID: cred.ID, UserID: h.ID, Name: name, Credential: pg.JSON[webauthn.Credential]{V: *cred}, Created: time.Now()}).Error
	})
	if errors.Is(err, errStale) {
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
	a.notify(a.passkeyMail(r.Context(), h.Email, name, true))
	w.WriteHeader(http.StatusNoContent)
}

func passkeyID(r *http.Request) ([]byte, bool) {
	id, err := base64.RawURLEncoding.DecodeString(r.PathValue("id"))
	return id, err == nil && len(id) > 0
}

func (a *Auth) rename(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("name"))
	s, ok := a.session(r)
	if !ok {
		a.toSignin(w, r)
		return
	}
	id, ok := passkeyID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if name == "" || len([]rune(name)) > 64 {
		a.accountPage(w, r, s, http.StatusUnprocessableEntity, "Passkey names need 1 to 64 characters.")
		return
	}
	res := a.db.WithContext(r.Context()).Model(&passkey{}).Where("id = ? AND user_id = ?", id, s.user.ID).Update("name", name)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected == 0 {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func (a *Auth) remove(w http.ResponseWriter, r *http.Request) {
	s, ok := a.session(r)
	if !ok {
		a.toSignin(w, r)
		return
	}
	if !s.recent() {
		a.accountPage(w, r, s, http.StatusForbidden, "Sign in again to remove a passkey.")
		return
	}
	id, ok := passkeyID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	var name string
	last := false
	err := a.db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT 1 FROM auth_users WHERE id = ? FOR NO KEY UPDATE", s.user.ID).Error; err != nil {
			return err
		}
		var k passkey
		if err := tx.Where("id = ? AND user_id = ?", id, s.user.ID).Take(&k).Error; err != nil {
			return err
		}
		if a.c.Mode == Operator {
			var n int64
			if err := tx.Model(&passkey{}).Where("user_id = ?", s.user.ID).Count(&n).Error; err != nil {
				return err
			}
			if n <= 1 {
				last = true
				return nil
			}
		}
		name = k.Name
		return tx.Delete(&passkey{}, "id = ? AND user_id = ?", id, s.user.ID).Error
	})
	if errors.Is(err, pg.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		panic(err)
	}
	if last {
		a.accountPage(w, r, s, http.StatusConflict, "This is your only way to sign in. Add another passkey before removing it.")
		return
	}
	a.notify(a.passkeyMail(r.Context(), s.user.Email, name, false))
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func encodeID(id []byte) string { return base64.RawURLEncoding.EncodeToString(id) }
