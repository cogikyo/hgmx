package auth

import (
	"errors"
	"net/http"
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

func (a *Auth) accountPage(w http.ResponseWriter, r *http.Request, s session, status int, problem string) {
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
	v := accountView{
		Script:   a.script,
		Email:    s.user.Email,
		Next:     next(r.URL.Query().Get("next")),
		Others:   int(others),
		Recent:   s.recent(),
		KeepLast: a.c.Mode == Operator && len(keys) == 1,
		Delete:   a.c.Mode == Open,
		Problem:  problem,
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
		// Recheck under the address lock: reset may have revoked this session.
		cur, ok, err := current(tx, r, s.user.ID)
		if err != nil {
			return err
		}
		if !ok || !cur.recent() {
			return errStale
		}
		if err := tx.Exec("UPDATE auth_codes SET consumed = ? WHERE email = ? AND consumed IS NULL", time.Now(), s.user.Email).Error; err != nil {
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
