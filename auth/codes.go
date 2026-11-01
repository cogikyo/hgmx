package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/pg"
)

const (
	signin = "signin"
	link   = "link"

	sendsPerMinute = 1
	sendsPerHour   = 5
	triesPerCode   = 5
	triesPerDay    = 10
)

type code struct {
	ID       pg.ID
	Email    string
	Purpose  string
	Mac      []byte
	Attempts int
	Created  time.Time
	Expires  time.Time
	Consumed *time.Time
}

func (code) TableName() string { return "auth_codes" }

func mac(key []byte, id pg.ID, purpose, email, secret string) []byte {
	m := hmac.New(sha256.New, key)
	fmt.Fprintf(m, "%s\x00%s\x00%s\x00%s", id, purpose, email, secret)
	return m.Sum(nil)
}

func newCode(tx *gorm.DB, key []byte, email, purpose, secret string, life time.Duration) (pg.ID, error) {
	now := time.Now()
	id := pg.NewID()
	c := code{ID: id, Email: email, Purpose: purpose, Mac: mac(key, id, purpose, email, secret), Created: now, Expires: now.Add(life)}
	return id, tx.Create(&c).Error
}

type sent struct {
	code    string
	limited bool
	deliver bool
	known   bool
}

func (a *Auth) issueCode(ctx context.Context, email string) (sent, error) {
	var s sent
	db := a.db.WithContext(ctx)
	if err := sweep(db, a.library); err != nil {
		return s, err
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, email); err != nil {
			return err
		}
		now := time.Now()
		var n struct{ Minute, Hour int }
		err := tx.Raw(`SELECT count(*) FILTER (WHERE created > ?) AS minute, count(*) AS hour
FROM auth_codes WHERE email = ? AND purpose = ? AND created > ?`,
			now.Add(-time.Minute), email, signin, now.Add(-time.Hour)).Scan(&n).Error
		if err != nil {
			return err
		}
		if n.Minute >= sendsPerMinute || n.Hour >= sendsPerHour {
			s.limited = true
			return nil
		}
		if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_users WHERE email = ?)", email).Scan(&s.known).Error; err != nil {
			return err
		}
		d, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
		if err != nil {
			return err
		}
		s.code = fmt.Sprintf("%06d", d.Int64())
		s.deliver = s.known || a.c.Mode == Open
		// Undelivered codes keep unknown addresses on the same rate-accounting path.
		_, err = newCode(tx, a.c.Key, email, signin, s.code, codeLife)
		return err
	})
	return s, err
}

type verdict int

const (
	accepted verdict = iota
	wrong
	stale
	exhausted
)

type verified struct {
	verdict  verdict
	token    string
	passkeys bool
}

func (a *Auth) checkCode(ctx context.Context, r *http.Request, email, guess string) (verified, error) {
	var v verified
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, email); err != nil {
			return err
		}
		now := time.Now()
		var tries int
		err := tx.Raw("SELECT coalesce(sum(attempts), 0) FROM auth_codes WHERE email = ? AND purpose = ? AND created > ?",
			email, signin, now.Add(-retention)).Scan(&tries).Error
		if err != nil {
			return err
		}
		if tries >= triesPerDay {
			v.verdict = exhausted
			return nil
		}
		var c code
		err = tx.Where("email = ? AND purpose = ? AND created > ?", email, signin, now.Add(-retention)).
			Order("created DESC, id DESC").Take(&c).Error
		if errors.Is(err, pg.ErrNotFound) {
			v.verdict = stale
			return nil
		}
		if err != nil {
			return err
		}
		if c.Consumed != nil || !now.Before(c.Expires) || c.Attempts >= triesPerCode {
			v.verdict = stale
			return nil
		}
		if !hmac.Equal(c.Mac, mac(a.c.Key, c.ID, signin, email, guess)) {
			v.verdict = wrong
			return tx.Exec("UPDATE auth_codes SET attempts = attempts + 1 WHERE id = ?", c.ID).Error
		}
		if err := tx.Exec("UPDATE auth_codes SET consumed = ? WHERE id = ?", now, c.ID).Error; err != nil {
			return err
		}
		u, err := findUser(tx, email)
		if errors.Is(err, pg.ErrNotFound) && a.c.Mode == Open {
			u, err = createUser(tx, email)
		}
		if errors.Is(err, pg.ErrNotFound) {
			v.verdict = stale
			return nil
		}
		if err != nil {
			return err
		}
		if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_passkeys WHERE user_id = ?)", u.ID).Scan(&v.passkeys).Error; err != nil {
			return err
		}
		v.token, err = startSession(tx, r, u.ID)
		return err
	})
	return v, err
}

func issueLink(ctx context.Context, db *gorm.DB, key []byte, email string, reset bool) (string, error) {
	secret := rand.Text()
	var id pg.ID
	db = db.WithContext(ctx)
	if err := sweep(db, library()); err != nil {
		return "", err
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, email); err != nil {
			return err
		}
		u, err := findUser(tx, email)
		if errors.Is(err, pg.ErrNotFound) {
			u, err = createUser(tx, email)
		}
		if err != nil {
			return err
		}
		if reset {
			for _, q := range []string{"DELETE FROM auth_passkeys WHERE user_id = ?", "DELETE FROM auth_sessions WHERE user_id = ?", "DELETE FROM auth_ceremonies WHERE user_id = ?"} {
				if err := tx.Exec(q, u.ID).Error; err != nil {
					return err
				}
			}
			if err := tx.Exec("UPDATE auth_codes SET consumed = ? WHERE email = ? AND consumed IS NULL", time.Now(), email).Error; err != nil {
				return err
			}
		}
		id, err = newCode(tx, key, email, link, secret, linkLife)
		return err
	})
	return id.String() + "." + secret, err
}

func (a *Auth) redeem(ctx context.Context, r *http.Request, raw string) (string, bool, error) {
	idText, secret, _ := strings.Cut(raw, ".")
	id, err := pg.ParseID(idText)
	if err != nil || secret == "" {
		return "", false, nil
	}
	db := a.db.WithContext(ctx)
	var email string
	if err := db.Raw("SELECT email FROM auth_codes WHERE id = ? AND purpose = ?", id, link).Scan(&email).Error; err != nil || email == "" {
		return "", false, err
	}
	var t string
	ok := false
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, email); err != nil {
			return err
		}
		var c code
		// Reset may have consumed the link while we waited for the address lock.
		err := tx.Raw("SELECT * FROM auth_codes WHERE id = ? FOR UPDATE", id).Take(&c).Error
		if errors.Is(err, pg.ErrNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		now := time.Now()
		if c.Consumed != nil || !now.Before(c.Expires) || c.Attempts >= triesPerCode {
			return nil
		}
		if !hmac.Equal(c.Mac, mac(a.c.Key, c.ID, link, c.Email, secret)) {
			return tx.Exec("UPDATE auth_codes SET attempts = attempts + 1 WHERE id = ?", c.ID).Error
		}
		if err := tx.Exec("UPDATE auth_codes SET consumed = ? WHERE id = ?", now, c.ID).Error; err != nil {
			return err
		}
		u, err := findUser(tx, c.Email)
		if errors.Is(err, pg.ErrNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		t, err = startSession(tx, r, u.ID)
		ok = err == nil
		return err
	})
	return t, ok, err
}

func (a *Auth) sendCode(w http.ResponseWriter, r *http.Request) {
	v := codeView{Next: next(r.FormValue("next")), Open: a.c.Mode == Open}
	email, ok := parseEmail(r.FormValue("email"))
	if !ok {
		a.page(w, r, http.StatusUnprocessableEntity, "Sign in", signinPage(a.signinView(r, v.Next, r.FormValue("email"), "Enter a valid email address.")))
		return
	}
	v.Email = email
	if !a.sends.allow(r) {
		a.page(w, r, http.StatusTooManyRequests, "Sign in", signinPage(a.signinView(r, v.Next, email, "Too many codes from this network. Try again later.")))
		return
	}
	s, err := a.issueCode(r.Context(), email)
	if err != nil {
		panic(err)
	}
	if s.limited {
		v.Problem = "A code was sent recently. Use the latest one, or ask again in a little while."
		a.page(w, r, http.StatusTooManyRequests, "Enter your code", codePage(v))
		return
	}
	m := a.codeMail(r.Context(), email, s.code)
	if s.deliver {
		a.mail.Enqueue(m)
	}
	a.page(w, r, http.StatusOK, "Enter your code", codePage(v))
}

func (a *Auth) verifyCode(w http.ResponseWriter, r *http.Request) {
	v := codeView{Next: next(r.FormValue("next")), Open: a.c.Mode == Open}
	email, ok := parseEmail(r.FormValue("email"))
	if !ok {
		a.page(w, r, http.StatusUnprocessableEntity, "Sign in", signinPage(a.signinView(r, v.Next, "", "Enter a valid email address.")))
		return
	}
	v.Email = email
	guess := strings.Join(strings.Fields(r.FormValue("code")), "")
	res, err := a.checkCode(r.Context(), r, email, guess)
	if err != nil {
		panic(err)
	}
	switch res.verdict {
	case accepted:
		if res.passkeys {
			a.notify(a.signedInMail(r.Context(), email))
		}
		setCookie(w, sessionCookie, res.token, sessionLife)
		http.Redirect(w, r, v.Next, http.StatusSeeOther)
		return
	case wrong:
		v.Problem = "That code isn't right. Check the latest email and try again."
		a.page(w, r, http.StatusUnprocessableEntity, "Enter your code", codePage(v))
	case stale:
		v.Problem = "That code has expired or was replaced. Ask for a new one."
		a.page(w, r, http.StatusUnprocessableEntity, "Enter your code", codePage(v))
	case exhausted:
		v.Problem = "Too many wrong codes for this address. Try again tomorrow."
		a.page(w, r, http.StatusTooManyRequests, "Enter your code", codePage(v))
	}
}

func (a *Auth) showLink(w http.ResponseWriter, r *http.Request) {
	a.page(w, r, http.StatusOK, "Sign in", linkPage(linkView{Script: a.script}))
}

func (a *Auth) redeemLink(w http.ResponseWriter, r *http.Request) {
	t, ok, err := a.redeem(r.Context(), r, r.FormValue("token"))
	if err != nil {
		panic(err)
	}
	if !ok {
		a.page(w, r, http.StatusNotFound, "Sign in", linkPage(linkView{Script: a.script, Problem: "This link has expired or was already used. Ask for a new one."}))
		return
	}
	setCookie(w, sessionCookie, t, sessionLife)
	http.Redirect(w, r, accountPath+"?"+url.Values{"next": {"/"}}.Encode(), http.StatusSeeOther)
}
