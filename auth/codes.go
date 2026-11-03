package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/clientip"
	"github.com/cogikyo/hgmx/pg"
)

const (
	signin = "signin"
	link   = "link"
	verify = "verify"

	sendsPerMinute = 1
	sendsPerHour   = 5
	triesPerCode   = 5
	triesPerDay    = 10
)

var digits = regexp.MustCompile(`^[0-9]{6}$`)

type code struct {
	ID        pg.ID
	Origin    string
	Email     string
	Purpose   string
	Mac       []byte
	Attempts  int
	Created   time.Time
	Expires   time.Time
	Consumed  *time.Time
	SessionID *pg.ID
	UserID    *pg.ID
}

func (code) TableName() string { return "auth_codes" }

func mac(key []byte, id pg.ID, purpose, email, secret string) []byte {
	m := hmac.New(sha256.New, key)
	fmt.Fprintf(m, "%s\x00%s\x00%s\x00%s", id, purpose, email, secret)
	return m.Sum(nil)
}

func mailSecret(key []byte, c code) string {
	return base64.RawURLEncoding.EncodeToString(mac(key, c.ID, "mail", c.subject(), ""))
}

func (c code) subject() string {
	s := c.Origin + "\x00" + c.Email
	for _, id := range []*pg.ID{c.SessionID, c.UserID} {
		if id != nil {
			s += "\x00" + id.String()
		}
	}
	return s
}

func newCode(tx *gorm.DB, key []byte, c code, secret string, life time.Duration) (code, error) {
	now := time.Now()
	c.ID, c.Created, c.Expires = pg.NewID(), now, now.Add(life)
	c.Mac = mac(key, c.ID, c.Purpose, c.subject(), secret)
	return c, tx.Create(&c).Error
}

func revokeCodes(tx *gorm.DB, u User) error {
	return tx.Exec("UPDATE auth_codes SET consumed = ? WHERE (user_id = ? OR email = ? AND email <> '') AND consumed IS NULL", time.Now(), u.ID, u.Email).Error
}

func throttled(tx *gorm.DB, origin, email string) (bool, error) {
	now := time.Now()
	var n struct{ Minute, Hour int }
	err := tx.Raw(`SELECT count(*) FILTER (WHERE created > ?) AS minute, count(*) AS hour
FROM auth_codes WHERE origin = ? AND email = ? AND purpose IN (?, ?) AND created > ?`,
		now.Add(-time.Minute), origin, email, signin, verify, now.Add(-time.Hour)).Scan(&n).Error
	return n.Minute >= sendsPerMinute || n.Hour >= sendsPerHour, err
}

type sent struct {
	code    string
	link    string
	limited bool
	deliver bool
	known   bool
}

func (a *Auth) issueCode(ctx context.Context, h *host, email string) (sent, error) {
	var s sent
	db := a.db.WithContext(ctx)
	if err := sweep(db, a.library); err != nil {
		return s, err
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, email); err != nil {
			return err
		}
		var err error
		if s.limited, err = throttled(tx, h.origin, email); err != nil || s.limited {
			return err
		}
		if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_users WHERE email = ?)", email).Scan(&s.known).Error; err != nil {
			return err
		}
		d, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
		if err != nil {
			return err
		}
		s.code = fmt.Sprintf("%06d", d.Int64())
		s.deliver = s.known || h.mode == Open
		// Undelivered codes keep unknown addresses on the same rate-accounting path.
		c, err := newCode(tx, a.key, code{Origin: h.origin, Email: email, Purpose: signin}, s.code, codeLife)
		s.link = c.ID.String() + "." + mailSecret(a.key, c)
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
	taken
	unknown
	linked
	denied
)

type verified struct {
	verdict  verdict
	user     User
	token    string
	passkeys bool
}

func enter(tx *gorm.DB, c caller, email string) (verified, error) {
	u, err := findUser(tx, email)
	if errors.Is(err, pg.ErrNotFound) && c.host.mode == Open {
		u, err = createUser(tx, email)
	}
	if errors.Is(err, pg.ErrNotFound) {
		return verified{verdict: stale}, nil
	}
	if err != nil {
		return verified{verdict: stale}, err
	}
	return admit(tx, c, u)
}

func admit(tx *gorm.DB, c caller, u User) (verified, error) {
	v := verified{verdict: stale}
	if !c.host.admits(u) {
		return verified{verdict: denied, user: u}, errDenied
	}
	if err := claim(tx, "id", u.ID); err != nil {
		return v, err
	}
	if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_passkeys WHERE user_id = ?)", u.ID).Scan(&v.passkeys).Error; err != nil {
		return v, err
	}
	var err error
	v.token, err = startSession(tx, c, u.ID)
	v.verdict, v.user = accepted, u
	return v, err
}

func (a *Auth) checkCode(ctx context.Context, k caller, email, guess string) (verified, error) {
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
		err = tx.Where("origin = ? AND email = ? AND purpose = ? AND created > ?", k.host.origin, email, signin, now.Add(-retention)).
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
		if !hmac.Equal(c.Mac, mac(a.key, c.ID, signin, c.subject(), guess)) {
			v.verdict = wrong
			return tx.Exec("UPDATE auth_codes SET attempts = attempts + 1 WHERE id = ?", c.ID).Error
		}
		if err := tx.Exec("UPDATE auth_codes SET consumed = ? WHERE id = ?", now, c.ID).Error; err != nil {
			return err
		}
		v, err = enter(tx, k, email)
		return err
	})
	return v, admitted(err)
}

func admitted(err error) error {
	if errors.Is(err, errDenied) {
		return nil
	}
	return err
}

var errMoved = errors.New("auth: the account's email changed; run the command again")

func issueLink(ctx context.Context, db *gorm.DB, key []byte, h *host, email string, want *pg.ID, reset bool) (string, User, error) {
	secret := rand.Text()
	var id pg.ID
	var u User
	db = db.WithContext(ctx)
	if err := sweep(db, library()); err != nil {
		return "", u, err
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, email); err != nil {
			return err
		}
		var err error
		if want != nil {
			u, err = hold(tx, *want)
			if errors.Is(err, pg.ErrNotFound) || err == nil && u.Email != email {
				return errMoved
			}
		} else {
			u, err = findUser(tx, email)
			if errors.Is(err, pg.ErrNotFound) {
				u, err = createUser(tx, email)
			}
			if err == nil {
				u, err = hold(tx, u.ID)
			}
		}
		if err != nil {
			return err
		}
		if !h.admits(u) {
			return errDenied
		}
		if reset {
			for _, q := range []string{"DELETE FROM auth_passkeys WHERE user_id = ?", flowsOf, "DELETE FROM auth_sessions WHERE user_id = ?", "DELETE FROM auth_ceremonies WHERE user_id = ?", "DELETE FROM auth_identities WHERE user_id = ?"} {
				if err := tx.Exec(q, u.ID).Error; err != nil {
					return err
				}
			}
			if err := revokeCodes(tx, u); err != nil {
				return err
			}
		}
		c, err := newCode(tx, key, code{Origin: h.origin, Email: u.Email, Purpose: link, UserID: &u.ID}, secret, linkLife)
		id = c.ID
		return err
	})
	return id.String() + "." + secret, u, err
}

func parseLink(raw string) (pg.ID, string, bool) {
	idText, secret, _ := strings.Cut(raw, ".")
	id, err := pg.ParseID(idText)
	return id, secret, err == nil && secret != ""
}

// spend consumes link id when secret opens it; callers hold its address lock.
func (a *Auth) spend(tx *gorm.DB, h *host, id pg.ID, secret string) (code, bool, error) {
	var c code
	// Reset may have consumed the link while we waited for the address lock.
	err := tx.Raw("SELECT * FROM auth_codes WHERE id = ? AND origin = ? FOR UPDATE", id, h.origin).Take(&c).Error
	if errors.Is(err, pg.ErrNotFound) {
		return c, false, nil
	}
	if err != nil {
		return c, false, err
	}
	var valid bool
	if c.Purpose == signin {
		valid = h.mode.codes() && hmac.Equal([]byte(secret), []byte(mailSecret(a.key, c)))
	} else {
		valid = hmac.Equal(c.Mac, mac(a.key, c.ID, c.Purpose, c.subject(), secret))
	}
	now := time.Now()
	if !valid || c.Consumed != nil || !now.Before(c.Expires) {
		return c, false, nil
	}
	return c, true, tx.Exec("UPDATE auth_codes SET consumed = ? WHERE id = ?", now, c.ID).Error
}

func (a *Auth) redeem(ctx context.Context, r *http.Request, h *host, raw string) (verified, bool, error) {
	v := verified{verdict: stale}
	id, secret, ok := parseLink(raw)
	if !ok {
		return v, false, nil
	}
	db := a.db.WithContext(ctx)
	var c code
	if err := db.Raw("SELECT email, purpose, user_id FROM auth_codes WHERE id = ? AND origin = ?", id, h.origin).Scan(&c).Error; err != nil {
		return v, false, err
	}
	mailed := c.Purpose == signin
	if !mailed && (c.Purpose != link || c.UserID == nil) {
		return v, false, nil
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, c.Email); err != nil {
			return err
		}
		var u User
		if !mailed {
			var err error
			if u, err = hold(tx, *c.UserID); err != nil {
				return err
			}
		}
		c, ok, err := a.spend(tx, h, id, secret)
		if err != nil || !ok {
			return err
		}
		if mailed {
			v, err = enter(tx, browser(r, h), c.Email)
		} else {
			v, err = admit(tx, browser(r, h), u)
		}
		return err
	})
	if errors.Is(err, pg.ErrNotFound) {
		return verified{verdict: stale}, mailed, nil
	}
	return v, mailed, admitted(err)
}

func (a *Auth) sendCode(w http.ResponseWriter, r *http.Request, h *host) {
	if !h.mode.codes() {
		http.NotFound(w, r)
		return
	}
	v := CodeView{Next: next(r.FormValue("next")), Open: h.mode == Open}
	email, ok := parseEmail(r.FormValue("email"))
	if !ok {
		a.signinPage(w, r, h, http.StatusUnprocessableEntity, v.Next, r.FormValue("email"), "Enter a valid email address.")
		return
	}
	v.Email = email
	if !h.sends.allow(clientip.From(r)) {
		a.signinPage(w, r, h, http.StatusTooManyRequests, v.Next, email, "Too many codes from this network. Try again later.")
		return
	}
	s, err := a.issueCode(r.Context(), h, email)
	if err != nil {
		panic(err)
	}
	if s.limited {
		v.Problem = "A code was sent recently. Use the latest one, or ask again in a little while."
		a.codePage(w, r, h, http.StatusTooManyRequests, v)
		return
	}
	m := a.codeMail(r.Context(), email, s.code, h.origin+"/auth/link?"+url.Values{"next": {v.Next}}.Encode()+"#"+s.link)
	if s.deliver {
		a.mail.Enqueue(m)
	}
	a.codePage(w, r, h, http.StatusOK, v)
}

func (a *Auth) verifyCode(w http.ResponseWriter, r *http.Request, h *host) {
	if !h.mode.codes() {
		http.NotFound(w, r)
		return
	}
	v := CodeView{Next: next(r.FormValue("next")), Open: h.mode == Open}
	email, ok := parseEmail(r.FormValue("email"))
	if !ok {
		a.signinPage(w, r, h, http.StatusUnprocessableEntity, v.Next, "", "Enter a valid email address.")
		return
	}
	v.Email = email
	guess := strings.Join(strings.Fields(r.FormValue("code")), "")
	if !digits.MatchString(guess) {
		v.Problem = "Enter the 6-digit code from the email."
		a.codePage(w, r, h, http.StatusUnprocessableEntity, v)
		return
	}
	if !h.verifies.allow(clientip.From(r)) {
		v.Problem = "Too many tries from this network. Try again later."
		a.codePage(w, r, h, http.StatusTooManyRequests, v)
		return
	}
	res, err := a.checkCode(r.Context(), browser(r, h), email, guess)
	if err != nil {
		panic(err)
	}
	switch res.verdict {
	case accepted:
		if res.passkeys {
			a.notify(a.signedInMail(r.Context(), h, email))
		}
		setCookie(w, sessionCookie, res.token, sessionLife)
		http.Redirect(w, r, v.Next, http.StatusSeeOther)
		return
	case denied:
		a.forbid(w, r, h, res.user)
	case wrong:
		v.Problem = "That code isn't right. Check the latest email and try again."
		a.codePage(w, r, h, http.StatusUnprocessableEntity, v)
	case stale:
		v.Problem = "That code has expired or was replaced. Ask for a new one."
		a.codePage(w, r, h, http.StatusUnprocessableEntity, v)
	case exhausted:
		v.Problem = "Too many wrong codes for this address. Use the link in the email, or try again tomorrow."
		a.codePage(w, r, h, http.StatusTooManyRequests, v)
	}
}

func (a *Auth) showLink(w http.ResponseWriter, r *http.Request, h *host) {
	a.linkPage(w, r, h, http.StatusOK, next(r.URL.Query().Get("next")), "")
}

func (a *Auth) redeemLink(w http.ResponseWriter, r *http.Request, h *host) {
	v, mailed, err := a.redeem(r.Context(), r, h, r.FormValue("token"))
	if err != nil {
		panic(err)
	}
	switch v.verdict {
	case accepted:
	case denied:
		a.forbid(w, r, h, v.user)
		return
	default:
		a.linkPage(w, r, h, http.StatusNotFound, next(r.FormValue("next")), "This link has expired or was already used. Ask for a new one.")
		return
	}
	to := accountPath + "?" + url.Values{"next": {"/"}}.Encode()
	if mailed {
		if v.passkeys {
			a.notify(a.signedInMail(r.Context(), h, v.user.Email))
		}
		to = next(r.FormValue("next"))
	}
	setCookie(w, sessionCookie, v.token, sessionLife)
	http.Redirect(w, r, to, http.StatusSeeOther)
}
