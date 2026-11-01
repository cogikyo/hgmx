// Package auth provides per-site passwordless accounts and revocable sessions.
// Sites own access policy; [Auth.Require] takes an allow check.
// Use one Auth per database: sessions are not scoped to an origin.
//
// # Site wiring
//
// Run [Track] before the site's migration track, never at server startup.
// Read the site's auth-key with credential.Read("auth-key").
// Its development fallback is HGMX_AUTH_KEY; the key needs at least 32 bytes.
// Start queue with mail.NewQueue and pass queue.Close to server.Run's after hooks.
// Keep gorm in store/ by adding Store methods:
//
//	func (s *Store) Auth(q *mail.Queue, c auth.Config) (*auth.Auth, error) {
//		return auth.New(s.db, q, c)
//	}
//	func (s *Store) AuthCommand(ctx context.Context, c auth.Config, args []string, out io.Writer) error {
//		return auth.Command(ctx, s.db, c, args, out)
//	}
//
// In site wiring, st, queue and key are ready and layout is the site's document:
//
//	var a *auth.Auth
//	a, err := st.Auth(queue, auth.Config{
//		Name: "Example", Origin: "https://example.com", Mode: auth.Open,
//		Key: []byte(key),
//		Layout: func(r *http.Request, t string, b templ.Component) templ.Component {
//			u, _ := a.User(r)
//			return layout(t, u, b)
//		},
//	})
//	if err != nil { return err }
//	a.Mount(mux)
//	mux.Handle("GET /members", a.Require(auth.Anyone, members))
//
// Declare a before st.Auth so Layout can read the user when a page is rendered.
// Layout must return the full document, including the body component.
// Mount under server.Site's cross-origin protection, or provide equivalent CSRF checks.
//
// # Modes and operator commands
//
// [Operator], the zero value, accepts only passkeys and CLI sign-in links.
// [Invite] also accepts email codes, but only for existing users.
// [Open] creates a user on their first verified code and allows account deletion.
// Pass a non-nil mail queue in Invite and Open modes.
//
// In the site's CLI, call st.AuthCommand(ctx, cfg, args, out).
// Pass the arguments after "auth"; [Command] does not need Config.Layout.
// "auth link <email>" creates a missing user and prints a 30-minute sign-in link.
// "auth reset <email>" also removes passkeys and sessions and invalidates codes.
// Treat the printed link as a secret; reset also invalidates outstanding links.
//
// # Styling and browser support
//
// Auth ships no CSS; the site styles these hooks:
//   - .auth[data-auth]: page roots (signin, code, link, account or forbidden).
//   - .auth-field: field labels, including the delete confirmation.
//   - .auth-primary, .auth-secondary and .auth-danger: action buttons.
//   - .auth-note: supporting text; .auth-problem: error messages.
//   - .auth-passkeys: the passkey list.
//
// Passkey buttons use data-auth-passkey="login" or "register" and data-next.
// Keep data-auth-link-token, data-auth-problem and the data-auth page root.
// Auth serves passkey.js and loads it on the pages that need it.
// The site's CSP must allow script-src 'self' and connect-src 'self'.
// Passkeys need Chrome 129, Safari 18.4 or Firefox 119 or later.
// Older browsers can use email codes in Invite and Open modes.
package auth

import (
	"context"
	_ "embed"
	"errors"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/mail"
	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/static"
)

// Mode decides how people sign in and who may create an account.
type Mode int

const (
	// Operator is the default: passkeys and CLI sign-in links only.
	Operator Mode = iota
	// Invite adds email codes for existing users; auth link creates users.
	Invite
	// Open creates users on their first verified code and allows account deletion.
	Open
)

func (m Mode) codes() bool { return m != Operator }

// Config describes one site's accounts.
type Config struct {
	// Name is the site name shown by passkey managers and in mail.
	Name string
	// Origin is the public HTTPS origin; HTTP is allowed only on loopback hosts.
	// Its hostname is the WebAuthn RP ID.
	Origin string
	// Mode is Operator unless set to Invite or Open.
	Mode Mode
	// Key is the HMAC key for codes and links, at least 32 bytes.
	Key []byte
	// Layout wraps auth's pages in the site's document.
	Layout func(r *http.Request, title string, body templ.Component) templ.Component
}

// User is a signed-in account.
type User struct {
	ID      pg.ID
	Email   string
	Created time.Time
}

// Auth serves sign-in, the account page and session lookups.
type Auth struct {
	db      *gorm.DB
	mail    *mail.Queue
	c       Config
	origin  string
	wa      *webauthn.WebAuthn
	assets  *static.Server
	script  string
	sends   *limiter
	begins  *limiter
	library string
}

const (
	prefix       = "/auth/"
	signinPath   = prefix + "signin"
	accountPath  = prefix + "account"
	assetsPath   = prefix + "static/"
	recent       = 10 * time.Minute
	ceremonyLife = 5 * time.Minute
	codeLife     = 10 * time.Minute
	linkLife     = 30 * time.Minute
	retention    = 25 * time.Hour
	idleLife     = 30 * 24 * time.Hour
	sessionLife  = 90 * 24 * time.Hour
	seenEvery    = time.Hour
)

//go:embed passkey.js
var passkeyJS []byte

// Anyone allows every signed-in user.
func Anyone(*http.Request, User) bool { return true }

// New checks c and prepares handlers; it does not migrate the database.
// A nil q is allowed only in Operator mode and disables mail notices.
// The caller owns q and must close it on shutdown.
func New(db *gorm.DB, q *mail.Queue, c Config) (*Auth, error) {
	origin, err := c.check()
	if err != nil {
		return nil, err
	}
	if c.Layout == nil {
		return nil, errors.New("auth: Config.Layout is required")
	}
	if c.Mode.codes() && q == nil {
		return nil, errors.New("auth: Invite and Open modes need a mail queue")
	}
	wa, err := webauthn.New(&webauthn.Config{
		RPID:                  origin.Hostname(),
		RPDisplayName:         c.Name,
		RPOrigins:             []string{origin.String()},
		AttestationPreference: protocol.PreferNoAttestation,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			ResidentKey:        protocol.ResidentKeyRequirementRequired,
			RequireResidentKey: protocol.ResidentKeyRequired(),
			UserVerification:   protocol.VerificationRequired,
		},
	})
	if err != nil {
		return nil, err
	}
	assets := static.Map(assetsPath, map[string][]byte{"passkey.js": passkeyJS})
	script, _ := assets.URL("passkey.js")
	return &Auth{
		db:      db,
		mail:    q,
		c:       c,
		origin:  origin.String(),
		wa:      wa,
		assets:  assets,
		script:  script,
		sends:   newLimiter(20, time.Hour),
		begins:  newLimiter(30, 10*time.Minute),
		library: library(),
	}, nil
}

func (c Config) check() (*url.URL, error) {
	if c.Name == "" {
		return nil, errors.New("auth: Config.Name is required")
	}
	u, err := url.Parse(c.Origin)
	if err != nil || u.Host == "" || u.User != nil || strings.TrimSuffix(u.Path, "/") != "" || u.RawQuery != "" || u.Fragment != "" ||
		!(u.Scheme == "https" || u.Scheme == "http" && loopback(u.Hostname())) {
		return nil, errors.New("auth: Config.Origin must be https://host, or http:// on localhost")
	}
	if len(c.Key) < 32 {
		return nil, errors.New("auth: Config.Key must be at least 32 bytes")
	}
	if c.Mode < Operator || c.Mode > Open {
		return nil, errors.New("auth: Config.Mode must be Operator, Invite or Open")
	}
	return &url.URL{Scheme: u.Scheme, Host: u.Host}, nil
}

func loopback(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func library() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, d := range info.Deps {
		if d.Path == "github.com/go-webauthn/webauthn" {
			if d.Replace != nil {
				return d.Replace.Path + "@" + d.Replace.Version
			}
			return d.Version
		}
	}
	return ""
}

// Mount registers the /auth/ pages and endpoints on mux.
func (a *Auth) Mount(mux *http.ServeMux) {
	mux.Handle("GET "+assetsPath, a.assets)
	mux.HandleFunc("GET "+signinPath, a.signin)
	if a.c.Mode.codes() {
		mux.HandleFunc("POST /auth/code", a.sendCode)
		mux.HandleFunc("POST /auth/code/verify", a.verifyCode)
	}
	mux.HandleFunc("POST /auth/passkey/begin", a.loginBegin)
	mux.HandleFunc("POST /auth/passkey/finish", a.loginFinish)
	mux.HandleFunc("GET /auth/link", a.showLink)
	mux.HandleFunc("POST /auth/link", a.redeemLink)
	mux.HandleFunc("POST /auth/signout", a.signout)
	mux.HandleFunc("GET "+accountPath, a.account)
	mux.HandleFunc("POST /auth/account/passkeys/begin", a.addBegin)
	mux.HandleFunc("POST /auth/account/passkeys/finish", a.addFinish)
	mux.HandleFunc("POST /auth/account/passkeys/{id}/rename", a.rename)
	mux.HandleFunc("POST /auth/account/passkeys/{id}/remove", a.remove)
	mux.HandleFunc("POST /auth/account/signout-others", a.signoutOthers)
	if a.c.Mode == Open {
		mux.HandleFunc("POST /auth/account/delete", a.deleteAccount)
	}
}

type sessionKey struct{}

// Require serves h to signed-in users that allow accepts; others sign in or get 403.
// Pass Anyone to admit every signed-in user.
func (a *Auth) Require(allow func(*http.Request, User) bool, h http.Handler) http.Handler {
	if allow == nil {
		panic("auth: Require needs an allow check; use auth.Anyone for any signed-in user")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, ok := a.session(r)
		if !ok {
			http.Redirect(w, r, signinPath+"?"+url.Values{"next": {r.URL.RequestURI()}}.Encode(), http.StatusSeeOther)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), sessionKey{}, s))
		if !allow(r, s.user) {
			a.page(w, r, http.StatusForbidden, "Not allowed", forbidden(s.user.Email))
			return
		}
		h.ServeHTTP(w, r)
	})
}

// User returns the signed-in user, reusing Require's lookup when present.
func (a *Auth) User(r *http.Request) (User, bool) {
	s, ok := a.session(r)
	return s.user, ok
}
