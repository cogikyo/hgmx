package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/clientip"
	"github.com/cogikyo/hgmx/pg"
)

// Google configures browser sign-in; native apps use [App.Google].
type Google struct {
	ClientID     string // OAuth 2.0 client ID for a Google web application.
	ClientSecret string // Required; supplied by the site.
}

const (
	googleKind     = "google"
	googleCallback = prefix + "google/callback"
	flowCookie     = "__Host-oidc"
	flowLife       = 10 * time.Minute
)

var googleIssuer = oidc.ProviderConfig{
	IssuerURL:  "https://accounts.google.com",
	AuthURL:    "https://accounts.google.com/o/oauth2/v2/auth",
	TokenURL:   "https://oauth2.googleapis.com/token",
	JWKSURL:    "https://www.googleapis.com/oauth2/v3/certs",
	Algorithms: []string{oidc.RS256},
}

type provider struct {
	kind     string
	client   *http.Client
	keys     *oidc.Provider
	oauth    oauth2.Config
	verifier *oidc.IDTokenVerifier
}

func discover(issuer oidc.ProviderConfig) (*http.Client, *oidc.Provider) {
	client := &http.Client{Timeout: 10 * time.Second}
	return client, issuer.NewProvider(oidc.ClientContext(context.Background(), client))
}

func newProvider(kind string, issuer oidc.ProviderConfig, id, secret string) *provider {
	client, p := discover(issuer)
	return &provider{
		kind:   kind,
		client: client,
		keys:   p,
		oauth: oauth2.Config{
			ClientID:     id,
			ClientSecret: secret,
			Endpoint:     p.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "email"},
		},
		verifier: p.Verifier(&oidc.Config{ClientID: id}),
	}
}

func (p *provider) config(origin string) oauth2.Config {
	c := p.oauth
	c.RedirectURL = origin + googleCallback
	return c
}

type flow struct {
	ID        []byte
	Origin    string
	Provider  string
	Intent    string
	Verifier  string
	Nonce     string
	SessionID *pg.ID
	Next      string
	Expires   time.Time
}

func (flow) TableName() string { return "auth_flows" }

type identity struct {
	subject  string
	email    string
	verified bool
}

func (a *Auth) startFlow(w http.ResponseWriter, r *http.Request, h *host, intent string, session *pg.ID, to string) {
	db := a.db.WithContext(r.Context())
	if err := sweep(db, a.library); err != nil {
		panic(err)
	}
	state := token()
	f := flow{ID: hash(state), Origin: h.origin, Provider: a.google.kind, Intent: intent, Verifier: oauth2.GenerateVerifier(), Nonce: token(), SessionID: session, Next: to, Expires: time.Now().Add(flowLife)}
	if err := db.Create(&f).Error; err != nil {
		panic(err)
	}
	opts := []oauth2.AuthCodeOption{oauth2.S256ChallengeOption(f.Verifier), oidc.Nonce(f.Nonce)}
	if intent == link {
		opts = append(opts, oauth2.SetAuthURLParam("prompt", "select_account"))
	}
	setCookie(w, flowCookie, state, flowLife)
	c := a.google.config(h.origin)
	a.continuePage(w, r, h, c.AuthCodeURL(state, opts...))
}

func (a *Auth) takeFlow(w http.ResponseWriter, r *http.Request, h *host) (flow, bool) {
	var f flow
	c, err := r.Cookie(flowCookie)
	if err != nil {
		return f, false
	}
	setCookie(w, flowCookie, "", -time.Second)
	state := r.URL.Query().Get("state")
	if subtle.ConstantTimeCompare([]byte(state), []byte(c.Value)) != 1 {
		return f, false
	}
	if err := a.db.WithContext(r.Context()).Raw("DELETE FROM auth_flows WHERE id = ? AND origin = ? RETURNING *", hash(state), h.origin).Scan(&f).Error; err != nil {
		panic(err)
	}
	return f, f.ID != nil && f.Provider == a.google.kind && time.Now().Before(f.Expires)
}

func (p *provider) identify(ctx context.Context, f flow, code string) (identity, error) {
	var id identity
	ctx = oidc.ClientContext(ctx, p.client)
	c := p.config(f.Origin)
	t, err := c.Exchange(ctx, code, oauth2.VerifierOption(f.Verifier))
	if err != nil {
		return id, err
	}
	raw, ok := t.Extra("id_token").(string)
	if !ok {
		return id, errors.New("auth: token response has no id_token")
	}
	tok, err := p.verifier.Verify(ctx, raw)
	if err != nil {
		return id, err
	}
	if subtle.ConstantTimeCompare([]byte(tok.Nonce), []byte(f.Nonce)) != 1 {
		return id, errors.New("auth: id_token nonce mismatch")
	}
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := tok.Claims(&claims); err != nil {
		return id, err
	}
	return identity{subject: tok.Subject, email: claims.Email, verified: claims.Verified}, nil
}

func owner(tx *gorm.DB, kind, subject string) (pg.ID, bool, error) {
	var id pg.ID
	res := tx.Raw("SELECT user_id FROM auth_identities WHERE kind = ? AND subject = ?", kind, subject).Scan(&id)
	return id, res.RowsAffected == 1, res.Error
}

func addIdentity(tx *gorm.DB, user pg.ID, kind, subject string) error {
	return tx.Exec("INSERT INTO auth_identities (user_id, kind, subject, created) VALUES (?, ?, ?, ?)", user, kind, subject, time.Now()).Error
}

func baseName(email string) string {
	local, _, _ := strings.Cut(email, "@")
	name := strings.Map(func(r rune) rune {
		if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || strings.ContainsRune("._-", r) {
			return r
		}
		return -1
	}, local)
	if name == "" {
		return "user"
	}
	return name[:min(len(name), 48)]
}

func createNamed(tx *gorm.DB, base string) (User, error) {
	handle := make([]byte, 16)
	rand.Read(handle)
	u := User{ID: pg.NewID(), Username: base, Created: time.Now()}
	for range 10 {
		res := tx.Exec("INSERT INTO auth_users (id, username, handle, created) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING", u.ID, u.Username, handle, u.Created)
		if res.Error != nil || res.RowsAffected == 1 {
			return u, res.Error
		}
		u.Username = base + "-" + strings.ToLower(rand.Text()[:6])
	}
	return u, errors.New("auth: no free username")
}

func (a *Auth) enterGoogle(ctx context.Context, c caller, id identity) (verified, error) {
	v := verified{verdict: stale}
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user, found, err := owner(tx, googleKind, id.subject)
		if err != nil {
			return err
		}
		if found {
			if v.user, err = hold(tx, user); err != nil {
				return err
			}
			if again, ok, err := owner(tx, googleKind, id.subject); err != nil || !ok || again != user {
				return err
			}
		} else {
			if c.host.mode != Open {
				v.verdict = unknown
				return nil
			}
			if email, ok := parseEmail(id.email); ok && id.verified {
				var used bool
				if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_users WHERE email = ?)", email).Scan(&used).Error; err != nil || used {
					v.verdict = taken
					return err
				}
			}
			if v.user, err = createNamed(tx, baseName(id.email)); err != nil {
				return err
			}
			if err := addIdentity(tx, v.user.ID, googleKind, id.subject); err != nil {
				return err
			}
		}
		if !c.host.admits(v.user) {
			v.verdict = denied
			return errDenied
		}
		v.token, err = startSession(tx, c, v.user.ID)
		v.verdict = accepted
		return err
	})
	if errors.Is(err, pg.ErrNotFound) || errors.Is(err, pg.ErrDuplicate) {
		return verified{verdict: stale}, nil
	}
	return v, admitted(err)
}

func (a *Auth) linkGoogle(ctx context.Context, c caller, s session, from *pg.ID, id identity) (verdict, User, error) {
	v := stale
	var u User
	if from == nil || s.id != *from {
		return v, u, nil
	}
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		if u, err = hold(tx, s.user.ID); err != nil {
			return err
		}
		cur, err := live(tx, c, u.ID)
		if err != nil || cur.id != s.id || !cur.recent() {
			return err
		}
		user, found, err := owner(tx, googleKind, id.subject)
		if err != nil {
			return err
		}
		if found {
			v = taken
			if user == u.ID {
				v = accepted
			}
			return nil
		}
		var has bool
		if err := tx.Raw("SELECT EXISTS (SELECT 1 FROM auth_identities WHERE user_id = ? AND kind = ?)", u.ID, googleKind).Scan(&has).Error; err != nil || has {
			v = linked
			return err
		}
		v = accepted
		return addIdentity(tx, u.ID, googleKind, id.subject)
	})
	switch {
	case errors.Is(err, pg.ErrNotFound), errors.Is(err, errStale):
		return stale, u, nil
	case errors.Is(err, errDenied):
		return denied, u, nil
	case errors.Is(err, pg.ErrDuplicate):
		return taken, u, nil
	}
	return v, u, err
}

func (a *Auth) googleSignin(w http.ResponseWriter, r *http.Request, h *host) {
	to := next(r.FormValue("next"))
	if !h.begins.allow(clientip.From(r)) {
		a.signinPage(w, r, h, http.StatusTooManyRequests, to, "", "Too many sign-in attempts. Try again in a few minutes.")
		return
	}
	a.startFlow(w, r, h, signin, nil, to)
}

func (a *Auth) googleLink(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	if !s.recent() {
		a.accountPage(w, r, h, s, http.StatusForbidden, "Sign in again to link Google.")
		return
	}
	if !h.begins.allow(clientip.From(r)) {
		a.accountPage(w, r, h, s, http.StatusTooManyRequests, "Too many attempts. Try again in a few minutes.")
		return
	}
	a.startFlow(w, r, h, link, &s.id, accountPath)
}

func (a *Auth) googleCallback(w http.ResponseWriter, r *http.Request, h *host) {
	f, ok := a.takeFlow(w, r, h)
	if !ok {
		a.signinPage(w, r, h, http.StatusBadRequest, accountPath, "", "Google sign-in took too long or started in another browser. Try again.")
		return
	}
	fail := func(status int, msg string) {
		if f.Intent == signin {
			a.signinPage(w, r, h, status, f.Next, "", msg)
			return
		}
		s, ok := a.signedIn(w, r, h)
		if !ok {
			return
		}
		a.accountPage(w, r, h, s, status, msg)
	}
	q := r.URL.Query()
	if q.Get("error") != "" {
		fail(http.StatusBadRequest, "Google sign-in was cancelled.")
		return
	}
	id, err := a.google.identify(r.Context(), f, q.Get("code"))
	if err != nil {
		slog.WarnContext(r.Context(), "auth: google sign-in failed", "err", err)
		fail(http.StatusBadGateway, "Google sign-in didn't work. Try again.")
		return
	}
	if f.Intent == link {
		v, u := stale, User{}
		if s, ok := a.session(r, h); ok {
			if v, u, err = a.linkGoogle(r.Context(), browser(r, h), s, f.SessionID, id); err != nil {
				panic(err)
			}
		}
		switch v {
		case accepted:
			a.notify(a.googleMail(r.Context(), h, u.Email, true))
			http.Redirect(w, r, accountPath, http.StatusSeeOther)
		case denied:
			a.forbid(w, r, h, u)
		case taken:
			fail(http.StatusConflict, "That Google account is linked to another account here.")
		case linked:
			fail(http.StatusConflict, "Your account already has a Google account. Remove it before linking another.")
		default:
			fail(http.StatusForbidden, "Your sign-in changed while linking Google. Sign in again.")
		}
		return
	}
	v, err := a.enterGoogle(r.Context(), browser(r, h), id)
	if err != nil {
		panic(err)
	}
	switch v.verdict {
	case accepted:
		setCookie(w, sessionCookie, v.token, sessionLife)
		http.Redirect(w, r, f.Next, http.StatusSeeOther)
	case denied:
		a.forbid(w, r, h, v.user)
	case unknown:
		fail(http.StatusForbidden, "No account here uses that Google account. Sign in another way, then link Google from your account page.")
	case taken:
		fail(http.StatusConflict, "An account here already uses "+id.email+". Sign in the way you usually do, then link Google from your account page.")
	default:
		fail(http.StatusConflict, "Google sign-in changed while we checked it. Try again.")
	}
}

func (a *Auth) googleUnlink(w http.ResponseWriter, r *http.Request, h *host) {
	s, ok := a.signedIn(w, r, h)
	if !ok {
		return
	}
	if !s.recent() {
		a.accountPage(w, r, h, s, http.StatusForbidden, "Sign in again to remove Google.")
		return
	}
	u, err := a.unlinkGoogle(r.Context(), browser(r, h), s.user.ID)
	last := errors.Is(err, errLast)
	if last {
		err = nil
	}
	if errors.Is(err, errMissing) {
		http.Redirect(w, r, accountPath, http.StatusSeeOther)
		return
	}
	if a.refused(w, r, h, s, err) {
		return
	}
	if last {
		a.accountPage(w, r, h, s, http.StatusConflict, "Google is your only way to sign in. Add a passkey before removing it.")
		return
	}
	a.notify(a.googleMail(r.Context(), h, u.Email, false))
	http.Redirect(w, r, accountPath, http.StatusSeeOther)
}

func (a *Auth) unlinkGoogle(ctx context.Context, c caller, user pg.ID) (User, error) {
	var u User
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		if u, err = hold(tx, user); err != nil {
			return err
		}
		if _, err := live(tx, c, u.ID); err != nil {
			return err
		}
		n, err := a.methods(tx, c.host, u)
		if err != nil {
			return err
		}
		if n <= 1 {
			return errLast
		}
		res := tx.Exec("DELETE FROM auth_identities WHERE user_id = ? AND kind = ?", u.ID, googleKind)
		if res.Error == nil && res.RowsAffected == 0 {
			return errMissing
		}
		return res.Error
	})
	return u, err
}
