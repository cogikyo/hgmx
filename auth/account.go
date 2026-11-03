package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/pg"
)

var errStale = errors.New("auth: session or ceremony no longer valid")

func (a *Auth) account(w http.ResponseWriter, r *http.Request) {
	s, ok := a.session(r)
	if !ok {
		a.toSignin(w, r)
		return
	}
	a.accountPage(w, r, s, http.StatusOK, "")
}

func (a *Auth) methods(db *gorm.DB, u User) (int64, error) {
	var n struct{ Passkeys, Identities int64 }
	err := db.Raw(`SELECT (SELECT count(*) FROM auth_passkeys WHERE user_id = ?) AS passkeys,
	(SELECT count(*) FROM auth_identities WHERE user_id = ? AND kind = ?) AS identities`, u.ID, u.ID, googleKind).Scan(&n).Error
	total := n.Passkeys
	if a.google != nil {
		total += n.Identities
	}
	if u.Email != "" && a.c.Mode.codes() {
		total++
	}
	return total, err
}

func (a *Auth) accountPage(w http.ResponseWriter, r *http.Request, s session, status int, problem string) {
	a.renderAccount(w, r, s, status, problem, "")
}

func (a *Auth) renderAccount(w http.ResponseWriter, r *http.Request, s session, status int, problem, notice string) {
	db := a.db.WithContext(r.Context())
	var keys []passkey
	if err := db.Select("id", "name", "created", "used").Where("user_id = ?", s.user.ID).Order("created").Find(&keys).Error; err != nil {
		panic(err)
	}
	var others int64
	now := time.Now()
	err := db.Table("auth_sessions").Where("user_id = ? AND id <> ? AND expires > ? AND seen > ?", s.user.ID, s.id, now, now.Add(-idleLife)).Count(&others).Error
	if err != nil {
		panic(err)
	}
	total, err := a.methods(db, s.user)
	if err != nil {
		panic(err)
	}
	var linked bool
	if a.google != nil {
		if err := db.Raw("SELECT EXISTS (SELECT 1 FROM auth_identities WHERE user_id = ? AND kind = ?)", s.user.ID, googleKind).Scan(&linked).Error; err != nil {
			panic(err)
		}
	}
	v := accountView{
		Script:     a.script,
		Name:       s.user.name(),
		Username:   s.user.Username,
		Email:      s.user.Email,
		Mail:       a.mail != nil,
		Next:       next(r.URL.Query().Get("next")),
		Others:     int(others),
		Recent:     s.recent(),
		KeepLast:   len(keys) == 1 && total == 1,
		Google:     a.google != nil,
		Linked:     linked,
		KeepGoogle: linked && total == 1,
		Delete:     a.c.Mode == Open,
		Problem:    problem,
		Notice:     notice,
	}
	for _, k := range keys {
		v.Passkeys = append(v.Passkeys, passkeyRow{ID: encodeID(k.ID), Name: k.Name, Created: k.Created, Used: k.Used})
	}
	a.page(w, r, status, "Account", accountPage(v))
}

func (a *Auth) deleteAccount(w http.ResponseWriter, r *http.Request) {
	confirmed := r.FormValue("confirm") == "yes"
	s, ok := a.session(r)
	if !ok {
		a.toSignin(w, r)
		return
	}
	if !s.recent() {
		a.accountPage(w, r, s, http.StatusForbidden, "Sign in again to delete your account.")
		return
	}
	if !confirmed {
		a.accountPage(w, r, s, http.StatusUnprocessableEntity, "Tick the box to confirm deleting your account.")
		return
	}
	err := a.db.WithContext(r.Context()).Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, s.user.Email); err != nil {
			return err
		}
		u, err := hold(tx, s.user.ID)
		if errors.Is(err, pg.ErrNotFound) {
			return errStale
		}
		if err != nil {
			return err
		}
		// Recheck under the locks: reset may have revoked this session.
		cur, ok, err := current(tx, r, s.user.ID)
		if err != nil {
			return err
		}
		if !ok || !cur.recent() || u.Email != s.user.Email {
			return errStale
		}
		if err := revokeCodes(tx, u); err != nil {
			return err
		}
		return tx.Exec("DELETE FROM auth_users WHERE id = ?", s.user.ID).Error
	})
	if errors.Is(err, errStale) {
		a.toSignin(w, r)
		return
	}
	if errors.Is(err, pg.ErrForeignKey) {
		a.accountPage(w, r, s, http.StatusConflict, "This account still owns data on this site, so it can't be deleted here.")
		return
	}
	if err != nil {
		panic(err)
	}
	setCookie(w, sessionCookie, "", -time.Second)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *Auth) setUsername(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("username"))
	s, ok := a.session(r)
	if !ok {
		a.toSignin(w, r)
		return
	}
	if !validName(name, s.user.Email) {
		msg := "Usernames need 1 to 64 letters, digits, dots, dashes or underscores."
		if strings.Contains(name, "@") {
			msg = "A username can contain @ only when it's your own email."
		}
		a.accountPage(w, r, s, http.StatusUnprocessableEntity, msg)
		return
	}
	res := a.db.WithContext(r.Context()).Exec("UPDATE auth_users SET username = ? WHERE id = ? AND coalesce(email, '') = ?", name, s.user.ID, s.user.Email)
	if errors.Is(res.Error, pg.ErrDuplicate) {
		a.accountPage(w, r, s, http.StatusConflict, "That username is taken.")
		return
	}
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected == 0 {
		a.accountPage(w, r, s, http.StatusConflict, "Your email changed while saving. Try again.")
		return
	}
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func (a *Auth) changeEmail(w http.ResponseWriter, r *http.Request) {
	raw := r.FormValue("email")
	s, ok := a.session(r)
	if !ok {
		a.toSignin(w, r)
		return
	}
	if !s.recent() {
		a.accountPage(w, r, s, http.StatusForbidden, "Sign in again to change your email.")
		return
	}
	to, ok := parseEmail(raw)
	if !ok {
		a.accountPage(w, r, s, http.StatusUnprocessableEntity, "Enter a valid email address.")
		return
	}
	if to == s.user.Email {
		a.accountPage(w, r, s, http.StatusUnprocessableEntity, "That's already your email.")
		return
	}
	if !a.sends.allow(r) {
		a.accountPage(w, r, s, http.StatusTooManyRequests, "Too many emails from this network. Try again later.")
		return
	}
	sent, err := a.issueVerify(r.Context(), r, s.user.ID, to)
	if errors.Is(err, errStale) {
		a.toSignin(w, r)
		return
	}
	if err != nil {
		panic(err)
	}
	if sent.limited {
		a.accountPage(w, r, s, http.StatusTooManyRequests, "An email was sent to that address recently. Use the latest link, or ask again in a little while.")
		return
	}
	if sent.deliver {
		a.notify(a.verifyMail(r.Context(), to, a.origin+confirmPath+"#"+sent.link))
	}
	a.renderAccount(w, r, s, http.StatusOK, "", "If no other account uses "+to+", we emailed it a link to confirm the change. It expires in 30 minutes.")
}

func (a *Auth) issueVerify(ctx context.Context, r *http.Request, user pg.ID, to string) (sent, error) {
	var s sent
	db := a.db.WithContext(ctx)
	if err := sweep(db, a.library); err != nil {
		return s, err
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, to); err != nil {
			return err
		}
		if _, err := hold(tx, user); err != nil {
			return err
		}
		cur, ok, err := current(tx, r, user)
		if err != nil {
			return err
		}
		if !ok || !cur.recent() {
			return errStale
		}
		if s.limited, err = throttled(tx, to); err != nil || s.limited {
			return err
		}
		var taken bool
		if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_users WHERE email = ?)", to).Scan(&taken).Error; err != nil || taken {
			return err
		}
		secret := rand.Text()
		id, err := newCode(tx, a.c.Key, to, verify, secret, linkLife, &cur.id, nil)
		s.link, s.deliver = id.String()+"."+secret, true
		return err
	})
	if errors.Is(err, pg.ErrNotFound) {
		err = errStale
	}
	return s, err
}

type change struct {
	verdict  verdict
	from, to string
}

func (a *Auth) confirm(ctx context.Context, raw string) (change, error) {
	ch := change{verdict: stale}
	id, secret, ok := parseLink(raw)
	if !ok {
		return ch, nil
	}
	db := a.db.WithContext(ctx)
	var c struct {
		Email     string
		SessionID pg.ID
		UserID    pg.ID
	}
	res := db.Raw(`SELECT c.email, c.session_id, s.user_id FROM auth_codes c JOIN auth_sessions s ON s.id = c.session_id
WHERE c.id = ? AND c.purpose = ?`, id, verify).Scan(&c)
	if res.Error != nil || res.RowsAffected == 0 {
		return ch, res.Error
	}
	var u User
	if err := db.Raw("SELECT "+userColumns+" FROM auth_users WHERE id = ?", c.UserID).Scan(&u).Error; err != nil || u.ID != c.UserID {
		return ch, err
	}
	addrs := []string{c.Email}
	if u.Email != "" {
		addrs = append(addrs, u.Email)
	}
	slices.Sort(addrs)
	addrs = slices.Compact(addrs)
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, e := range addrs {
			if err := lock(tx, e); err != nil {
				return err
			}
		}
		cur, err := hold(tx, u.ID)
		if errors.Is(err, pg.ErrNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		if cur.Email != u.Email {
			return nil
		}
		c, ok, err := a.spend(tx, id, secret)
		if err != nil || !ok {
			return err
		}
		now := time.Now()
		var live bool
		err = tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_sessions WHERE id = ? AND user_id = ? AND expires > ? AND seen > ?)",
			c.SessionID, u.ID, now, now.Add(-idleLife)).Scan(&live).Error
		if err != nil || !live {
			return err
		}
		var used bool
		if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_users WHERE email = ? AND id <> ?)", c.Email, u.ID).Scan(&used).Error; err != nil {
			return err
		}
		if used {
			ch.verdict = taken
			return nil
		}
		if err := tx.Exec("UPDATE auth_codes SET consumed = ? WHERE email IN ? AND consumed IS NULL", now, addrs).Error; err != nil {
			return err
		}
		name := cur.Username
		if name == "" || ownEmail(name, cur.Email) {
			name = c.Email
		}
		if err := tx.Exec("UPDATE auth_users SET email = ?, username = ? WHERE id = ?", c.Email, name, u.ID).Error; err != nil {
			return err
		}
		ch = change{verdict: accepted, from: cur.Email, to: c.Email}
		return nil
	})
	return ch, err
}

func (a *Auth) showConfirm(w http.ResponseWriter, r *http.Request) {
	a.page(w, r, http.StatusOK, "Confirm email", confirmPage(linkView{Script: a.script}))
}

func (a *Auth) confirmEmail(w http.ResponseWriter, r *http.Request) {
	ch, err := a.confirm(r.Context(), r.FormValue("token"))
	if err != nil {
		panic(err)
	}
	switch ch.verdict {
	case accepted:
		a.notify(a.changedMail(r.Context(), ch.from, ch.to))
		http.Redirect(w, r, accountPath, http.StatusSeeOther)
	case taken:
		a.page(w, r, http.StatusConflict, "Confirm email", confirmPage(linkView{Script: a.script, Problem: "Another account now uses that email."}))
	default:
		a.page(w, r, http.StatusNotFound, "Confirm email", confirmPage(linkView{Script: a.script, Problem: "This link has expired or was already used. Ask for a new one on your account page."}))
	}
}
