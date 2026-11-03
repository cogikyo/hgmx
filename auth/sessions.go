package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/pg"
)

const (
	sessionCookie  = "__Host-session"
	ceremonyCookie = "__Host-ceremony"
)

type session struct {
	id      pg.ID
	created time.Time
	user    User
}

func (s session) recent() bool { return time.Since(s.created) <= recent }

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

func lookup(db *gorm.DB, origin, token string) (session, time.Time, error) {
	now := time.Now()
	var row struct {
		ID          pg.ID
		Created     time.Time
		Seen        time.Time
		UserID      pg.ID
		Email       string
		Username    string
		UserCreated time.Time
	}
	err := db.Raw(`
SELECT s.id, s.created, s.seen, u.id AS user_id, coalesce(u.email, '') AS email, coalesce(u.username, '') AS username, u.created AS user_created
FROM auth_sessions s JOIN auth_users u ON u.id = s.user_id
WHERE s.token_hash = ? AND s.origin = ? AND s.expires > ? AND s.seen > ?`,
		hash(token), origin, now, now.Add(-idleLife)).Take(&row).Error
	s := session{id: row.ID, created: row.Created, user: User{ID: row.UserID, Email: row.Email, Username: row.Username, Created: row.UserCreated}}
	return s, row.Seen, err
}

func (a *Auth) session(r *http.Request, h *host) (session, bool) {
	if c, ok := r.Context().Value(sessionKey{}).(cached); ok && c.host == h {
		return c.session, true
	}
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return session{}, false
	}
	db := a.db.WithContext(r.Context())
	s, seen, err := lookup(db, h.origin, c.Value)
	if errors.Is(err, pg.ErrNotFound) {
		return session{}, false
	}
	if err != nil {
		panic(err)
	}
	if now := time.Now(); now.Sub(seen) >= seenEvery {
		if err := db.Exec("UPDATE auth_sessions SET seen = ? WHERE id = ? AND seen < ?", now, s.id, now.Add(-seenEvery)).Error; err != nil {
			panic(err)
		}
	}
	return s, true
}

func current(tx *gorm.DB, r *http.Request, h *host, user pg.ID) (session, bool, error) {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return session{}, false, nil
	}
	s, _, err := lookup(tx, h.origin, c.Value)
	if errors.Is(err, pg.ErrNotFound) {
		return session{}, false, nil
	}
	return s, err == nil && s.user.ID == user, err
}

func live(tx *gorm.DB, r *http.Request, h *host, user pg.ID) (session, error) {
	s, ok, err := current(tx, r, h, user)
	switch {
	case err != nil:
		return s, err
	case !ok:
		return s, errStale
	case !h.admits(s.user):
		return s, errDenied
	}
	return s, nil
}

func startSession(tx *gorm.DB, r *http.Request, h *host, user pg.ID) (string, error) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		if err := tx.Exec("DELETE FROM auth_sessions WHERE token_hash = ? AND origin = ?", hash(c.Value), h.origin).Error; err != nil {
			return "", err
		}
	}
	now := time.Now()
	t := token()
	err := tx.Exec("INSERT INTO auth_sessions (id, user_id, token_hash, origin, created, seen, expires) VALUES (?, ?, ?, ?, ?, ?, ?)",
		pg.NewID(), user, hash(t), h.origin, now, now, now.Add(sessionLife)).Error
	return t, err
}

func (a *Auth) signout(w http.ResponseWriter, r *http.Request, h *host) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		if err := a.db.WithContext(r.Context()).Exec("DELETE FROM auth_sessions WHERE token_hash = ? AND origin = ?", hash(c.Value), h.origin).Error; err != nil {
			panic(err)
		}
	}
	setCookie(w, sessionCookie, "", -time.Second)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *Auth) signoutOthers(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	err := a.db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
		if _, err := hold(tx, s.user.ID); err != nil {
			return err
		}
		cur, err := live(tx, r, h, s.user.ID)
		if err != nil {
			return err
		}
		return tx.Exec("DELETE FROM auth_sessions WHERE user_id = ? AND origin = ? AND id <> ?", cur.user.ID, h.origin, cur.id).Error
	})
	if a.refused(w, r, h, s, err) {
		return
	}
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}
