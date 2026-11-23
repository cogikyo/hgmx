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

	"github.com/cogikyo/hgmx/clientip"
	"github.com/cogikyo/hgmx/pg"
	authv1 "github.com/cogikyo/hgmx/proto/hgmx/auth/v1"
)

func (a *Auth) account(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	a.accountPage(w, r, h, s, http.StatusOK, "")
}

func (a *Auth) methods(db *gorm.DB, h *host, u User) (int64, error) {
	var keys []passkey
	if err := db.Select("credential").Where("user_id = ?", u.ID).Find(&keys).Error; err != nil {
		return 0, err
	}
	n, err := identities(db, u.ID)
	return a.tally(h, u, keys, n), err
}

func identities(db *gorm.DB, user pg.ID) (int64, error) {
	var n int64
	err := db.Raw("SELECT count(*) FROM auth_identities WHERE user_id = ? AND kind = ?", user, googleKind).Scan(&n).Error
	return n, err
}

func (a *Auth) tally(h *host, u User, keys []passkey, identities int64) int64 {
	var total int64
	for _, k := range keys {
		if compatible(k.Credential.V, a.rp) {
			total++
		}
	}
	if a.googleOn(h) {
		total += identities
	}
	if u.Email != "" && h.mode.codes() {
		total++
	}
	return total
}

type profile struct {
	user     User
	passkeys []passkey
	linked   bool
	methods  int64
}

func (a *Auth) loadProfile(db *gorm.DB, h *host, u User) (profile, error) {
	p := profile{user: u}
	if err := db.Select("id", "name", "credential", "created", "used").Where("user_id = ?", u.ID).Order("created").Find(&p.passkeys).Error; err != nil {
		return p, err
	}
	n, err := identities(db, u.ID)
	p.linked, p.methods = n > 0, a.tally(h, u, p.passkeys, n)
	return p, err
}

func (a *Auth) accountPage(w http.ResponseWriter, r *http.Request, h *host, s session, status int, problem string) {
	a.renderAccount(w, r, h, s, status, problem, "")
}

func (a *Auth) renderAccount(w http.ResponseWriter, r *http.Request, h *host, s session, status int, problem, notice string) {
	db := a.db.WithContext(r.Context())
	p, err := a.loadProfile(db, h, s.user)
	if err != nil {
		panic(err)
	}
	rows, err := a.sessions(db, h, s.user.ID)
	if err != nil {
		panic(err)
	}
	linked := a.google != nil && p.linked
	v := AccountView{
		Name:       p.user.name(),
		Username:   p.user.Username,
		Email:      p.user.Email,
		Mail:       a.mail != nil,
		Next:       next(r.URL.Query().Get("next")),
		Recent:     s.recent(),
		Google:     a.google != nil,
		Linked:     linked,
		KeepGoogle: linked && p.methods == 1,
		Delete:     h.mode == Open,
		Problem:    problem,
		Notice:     notice,
	}
	for _, row := range rows {
		v.Sessions = append(v.Sessions, SessionRow{ID: row.ID.String(), Device: a.device(row), Created: row.Created, Seen: row.Seen, Current: row.ID == s.id})
	}
	for _, k := range p.passkeys {
		row := PasskeyRow{ID: encodeID(k.ID), Name: k.Name, Created: k.Created, Used: k.Used}
		if compatible(k.Credential.V, a.rp) {
			v.Passkeys = append(v.Passkeys, row)
		} else {
			v.Older = append(v.Older, row)
		}
	}
	v.KeepLast = len(v.Passkeys) == 1 && p.methods == 1
	a.page(w, r, h, status, "Account", root("account", a.script, h.views.Account(v)))
}

func (a *Auth) device(row sessionTable) string {
	switch {
	case row.Label != nil:
		return *row.Label
	case row.ClientID == nil:
		return "Browser"
	case a.apps[*row.ClientID].kind == authv1.ClientKind_CLIENT_KIND_IOS:
		return "iPhone or iPad app"
	}
	return "Android app"
}

func (a *Auth) revokeSession(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	id, err := pg.ParseID(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = a.revoke(r.Context(), browser(r, h), s.user.ID, &id)
	if !errors.Is(err, errMissing) && a.refused(w, r, h, s, err) {
		return
	}
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func (a *Auth) deleteAccount(w http.ResponseWriter, r *http.Request, h *host) {
	if h.mode != Open {
		http.NotFound(w, r)
		return
	}
	confirmed := r.FormValue("confirm") == "yes"
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	var err error
	switch {
	case !s.recent():
		err = errRecent
	case !confirmed:
		err = errUnconfirmed
	default:
		err = a.erase(r.Context(), browser(r, h), s)
	}
	if a.refused(w, r, h, s, err) {
		return
	}
	setCookie(w, sessionCookie, "", -time.Second)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

const flowsOf = "DELETE FROM auth_native_flows WHERE session_id IN (SELECT id FROM auth_sessions WHERE user_id = ?)"

func (a *Auth) erase(ctx context.Context, c caller, s session) error {
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, s.user.Email); err != nil {
			return err
		}
		u, err := hold(tx, s.user.ID)
		if err != nil {
			return err
		}
		// Recheck under the locks: reset may have revoked this session.
		cur, err := live(tx, c, s.user.ID)
		if err != nil {
			return err
		}
		if !cur.recent() || u.Email != s.user.Email {
			return errRecent
		}
		if err := revokeCodes(tx, u); err != nil {
			return err
		}
		if err := tx.Exec(flowsOf, u.ID).Error; err != nil {
			return err
		}
		if a.erasure != nil {
			if err := a.erasure(tx, u.ID); err != nil {
				return err
			}
		}
		return tx.Exec("DELETE FROM auth_users WHERE id = ?", u.ID).Error
	})
	if errors.Is(err, pg.ErrForeignKey) {
		return errOwned
	}
	return err
}

func (a *Auth) setUsername(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	if a.refused(w, r, h, s, a.saveUsername(r.Context(), browser(r, h), s, r.FormValue("username"))) {
		return
	}
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func (a *Auth) saveUsername(ctx context.Context, c caller, s session, raw string) error {
	name := strings.TrimSpace(raw)
	if !validName(name, s.user.Email) {
		if strings.Contains(name, "@") {
			return errNameAt
		}
		return errName
	}
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := hold(tx, s.user.ID); err != nil {
			return err
		}
		cur, err := live(tx, c, s.user.ID)
		if err != nil {
			return err
		}
		if cur.user.Email != s.user.Email {
			return errChanged
		}
		u := cur.user
		u.Username = name
		if !c.host.admits(u) {
			return errDenied
		}
		return tx.Exec("UPDATE auth_users SET username = ? WHERE id = ?", name, u.ID).Error
	})
	if errors.Is(err, pg.ErrDuplicate) {
		return errNameTaken
	}
	return err
}

func (a *Auth) changeEmail(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	to, err := newEmail(r.FormValue("email"), s.user.Email)
	switch {
	case !s.recent():
		err = errRecent
	case err == nil && !h.sends.allow(clientip.From(r)):
		err = errThrottled
	case err == nil:
		err = a.issueVerify(r.Context(), browser(r, h), s.user.ID, to)
	}
	if a.refused(w, r, h, s, err) {
		return
	}
	a.renderAccount(w, r, h, s, http.StatusOK, "", "If no other account uses "+to+", we emailed it a link to confirm the change. It expires in 30 minutes.")
}

func newEmail(raw, current string) (string, error) {
	to, ok := parseEmail(raw)
	switch {
	case !ok:
		return "", errBadEmail
	case to == current:
		return "", errSameEmail
	}
	return to, nil
}

func (a *Auth) issueVerify(ctx context.Context, k caller, user pg.ID, to string) error {
	h := k.host
	var link string
	var taken bool
	db := a.db.WithContext(ctx)
	if err := sweep(db, a.library); err != nil {
		return err
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := lock(tx, to); err != nil {
			return err
		}
		if _, err := hold(tx, user); err != nil {
			return err
		}
		cur, err := live(tx, k, user)
		if err != nil {
			return err
		}
		if !cur.recent() {
			return errRecent
		}
		limited, err := throttled(tx, h.origin, to)
		if err != nil {
			return err
		}
		if limited {
			return errTooSoon
		}
		if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_users WHERE email = ?)", to).Scan(&taken).Error; err != nil {
			return err
		}
		secret := rand.Text()
		c, err := newCode(tx, a.key, code{Origin: h.origin, Email: to, Purpose: verify, SessionID: &cur.id}, secret, linkLife)
		link = c.ID.String() + "." + secret
		return err
	})
	if err != nil {
		return err
	}
	// Compose unconditionally so timing does not reveal whether the address is taken.
	m := a.verifyMail(ctx, to, h.origin+confirmPath+"#"+link)
	if !taken {
		a.notify(m)
	}
	return nil
}

func (a *Auth) confirm(ctx context.Context, h *host, raw string) (verdict, error) {
	v := stale
	var from, to string
	id, secret, ok := parseLink(raw)
	if !ok {
		return v, nil
	}
	db := a.db.WithContext(ctx)
	var c struct {
		Email     string
		SessionID pg.ID
		UserID    pg.ID
	}
	res := db.Raw(`SELECT c.email, c.session_id, s.user_id FROM auth_codes c JOIN auth_sessions s ON s.id = c.session_id
WHERE c.id = ? AND c.purpose = ? AND c.origin = ?`, id, verify, h.origin).Scan(&c)
	if res.Error != nil || res.RowsAffected == 0 {
		return v, res.Error
	}
	var u User
	if err := db.Raw("SELECT "+userColumns+" FROM auth_users WHERE id = ?", c.UserID).Scan(&u).Error; err != nil || u.ID != c.UserID {
		return v, err
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
		c, ok, err := a.spend(tx, h, id, secret)
		if err != nil || !ok {
			return err
		}
		now := time.Now()
		var active bool
		err = tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_sessions WHERE id = ? AND user_id = ? AND origin = ? AND expires > ? AND seen > ?)",
			c.SessionID, u.ID, h.origin, now, now.Add(-idleLife)).Scan(&active).Error
		if err != nil || !active {
			return err
		}
		var used bool
		if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_users WHERE email = ? AND id <> ?)", c.Email, u.ID).Scan(&used).Error; err != nil {
			return err
		}
		if used {
			v = taken
			return nil
		}
		want := cur
		want.Email = c.Email
		if want.Username == "" || ownEmail(want.Username, cur.Email) {
			want.Username = c.Email
		}
		if !h.admits(cur) || !h.admits(want) {
			v = denied
			return nil
		}
		if err := tx.Exec("UPDATE auth_codes SET consumed = ? WHERE email IN ? AND consumed IS NULL", now, addrs).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE auth_users SET email = ?, username = ? WHERE id = ?", want.Email, want.Username, u.ID).Error; err != nil {
			return err
		}
		v, from, to = accepted, cur.Email, want.Email
		return nil
	})
	if err == nil && v == accepted {
		a.notify(a.changedMail(ctx, h, from, to))
	}
	return v, err
}

func (a *Auth) showConfirm(w http.ResponseWriter, r *http.Request, h *host) {
	a.emailPage(w, r, h, http.StatusOK, "")
}

func (a *Auth) confirmEmail(w http.ResponseWriter, r *http.Request, h *host) {
	v, err := a.confirm(r.Context(), h, r.FormValue("token"))
	if err != nil {
		panic(err)
	}
	switch v {
	case accepted:
		http.Redirect(w, r, accountPath, http.StatusSeeOther)
	case taken:
		a.emailPage(w, r, h, http.StatusConflict, "Another account now uses that email.")
	case denied:
		a.emailPage(w, r, h, http.StatusForbidden, "This account can't use that email here.")
	default:
		a.emailPage(w, r, h, http.StatusNotFound, "This link has expired or was already used. Ask for a new one on your account page.")
	}
}
