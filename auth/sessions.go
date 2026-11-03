package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/clientip"
	"github.com/cogikyo/hgmx/pg"
)

const (
	sessionCookie  = "__Host-session"
	ceremonyCookie = "__Host-ceremony"
)

type session struct {
	id      pg.ID
	created time.Time
	expires time.Time
	user    User
}

func (s session) recent() bool { return time.Since(s.created) <= recent }

type caller struct {
	host  *host
	app   *app
	token string
	label *string
	ip    string
}

func browser(r *http.Request, h *host) caller {
	c := caller{host: h, ip: clientip.From(r)}
	if k, err := r.Cookie(sessionCookie); err == nil {
		c.token = k.Value
	}
	return c
}

func (c caller) client() *string {
	if c.app == nil {
		return nil
	}
	return &c.app.id
}

func token() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func hash(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

func setCookie(w http.ResponseWriter, name, value string, life time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(life / time.Second),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func lookup(db *gorm.DB, c caller) (session, time.Time, error) {
	if c.token == "" {
		return session{}, time.Time{}, pg.ErrNotFound
	}
	now := time.Now()
	var row struct {
		ID          pg.ID
		Created     time.Time
		Seen        time.Time
		Expires     time.Time
		UserID      pg.ID
		Email       string
		Username    string
		UserCreated time.Time
	}
	err := db.Raw(`
SELECT s.id, s.created, s.seen, s.expires, u.id AS user_id, coalesce(u.email, '') AS email, coalesce(u.username, '') AS username, u.created AS user_created
FROM auth_sessions s JOIN auth_users u ON u.id = s.user_id
WHERE s.token_hash = ? AND s.origin = ? AND s.client_id IS NOT DISTINCT FROM ? AND s.expires > ? AND s.seen > ?`,
		hash(c.token), c.host.origin, c.client(), now, now.Add(-idleLife)).Take(&row).Error
	s := session{id: row.ID, created: row.Created, expires: row.Expires, user: User{ID: row.UserID, Email: row.Email, Username: row.Username, Created: row.UserCreated}}
	return s, row.Seen, err
}

func (a *Auth) find(ctx context.Context, c caller) (session, bool, error) {
	db := a.db.WithContext(ctx)
	s, seen, err := lookup(db, c)
	if errors.Is(err, pg.ErrNotFound) {
		return session{}, false, nil
	}
	if err != nil {
		return session{}, false, err
	}
	if now := time.Now(); now.Sub(seen) >= seenEvery {
		if err := db.Exec("UPDATE auth_sessions SET seen = ? WHERE id = ? AND seen < ?", now, s.id, now.Add(-seenEvery)).Error; err != nil {
			return session{}, false, err
		}
	}
	return s, true, nil
}

func (a *Auth) session(r *http.Request, h *host) (session, bool) {
	if c, ok := r.Context().Value(sessionKey{}).(cached); ok && c.host == h {
		return c.session, true
	}
	s, ok, err := a.find(r.Context(), browser(r, h))
	if err != nil {
		panic(err)
	}
	return s, ok
}

func current(tx *gorm.DB, c caller, user pg.ID) (session, bool, error) {
	s, _, err := lookup(tx, c)
	if errors.Is(err, pg.ErrNotFound) {
		return session{}, false, nil
	}
	return s, err == nil && s.user.ID == user, err
}

func live(tx *gorm.DB, c caller, user pg.ID) (session, error) {
	s, ok, err := current(tx, c, user)
	switch {
	case err != nil:
		return s, err
	case !ok:
		return s, errStale
	case !c.host.admits(s.user):
		return s, errDenied
	}
	return s, nil
}

func startSession(tx *gorm.DB, c caller, user pg.ID) (string, error) {
	if err := end(tx, c); err != nil {
		return "", err
	}
	now := time.Now()
	t := token()
	err := tx.Exec("INSERT INTO auth_sessions (id, user_id, token_hash, origin, client_id, label, created, seen, expires) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		pg.NewID(), user, hash(t), c.host.origin, c.client(), c.label, now, now, now.Add(sessionLife)).Error
	return t, err
}

func end(db *gorm.DB, c caller) error {
	if c.token == "" {
		return nil
	}
	return db.Exec("DELETE FROM auth_sessions WHERE token_hash = ? AND origin = ? AND client_id IS NOT DISTINCT FROM ?", hash(c.token), c.host.origin, c.client()).Error
}

func (a *Auth) signout(w http.ResponseWriter, r *http.Request, h *host) {
	if err := end(a.db.WithContext(r.Context()), browser(r, h)); err != nil {
		panic(err)
	}
	setCookie(w, sessionCookie, "", -time.Second)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *Auth) signoutOthers(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	err := a.revoke(r.Context(), browser(r, h), s.user.ID, nil)
	if a.refused(w, r, h, s, err) {
		return
	}
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func (a *Auth) revoke(ctx context.Context, c caller, user pg.ID, id *pg.ID) error {
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := hold(tx, user); err != nil {
			return err
		}
		cur, err := live(tx, c, user)
		if err != nil {
			return err
		}
		if id == nil {
			return tx.Exec("DELETE FROM auth_sessions WHERE user_id = ? AND origin = ? AND id <> ?", user, c.host.origin, cur.id).Error
		}
		res := tx.Exec("DELETE FROM auth_sessions WHERE user_id = ? AND origin = ? AND id = ?", user, c.host.origin, *id)
		if res.Error == nil && res.RowsAffected == 0 {
			return errMissing
		}
		return res.Error
	})
}
