// Package auth provides passwordless accounts and origin-bound sessions.
// Use one Auth for all hosts that share a database.
// [Host.Allow] gates each host; [Auth.Require] adds route-level checks.
//
// # Site wiring
//
// Run [Track] before the site's migration track, never at server startup.
// Read the site's auth-key through a handle declared from [KeyCredential].
// Its fallback is HGMX_AUTH_KEY; the key needs at least 32 bytes.
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
//		Name: "Example", Key: []byte(key),
//		Hosts: []auth.Host{{
//			Origin: "https://example.com", Mode: auth.Open,
//			Layout: func(r *http.Request, t string, b templ.Component) templ.Component {
//				u, _ := a.User(r)
//				return layout(t, u, b)
//			},
//		}},
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
// [Invite] also accepts emailed codes and links, but only for existing users.
// [Open] creates a user on their first verified code or link and allows account deletion.
// Pass a non-nil mail queue in Invite and Open modes.
// Config.Google adds Google sign-in in every mode.
// Register every Host.Origin + "/auth/google/callback" with Google.
// Signed-in users can link Google; Open mode can create Google-only users.
// Google's email never becomes the user's email.
//
// In the site's CLI, call st.AuthCommand(ctx, cfg, args, out).
// Pass the arguments after "auth"; [Command] does not need Host.Layout.
// Use "auth link [--origin URL] <who>" or "auth reset [--origin URL] <who>".
// --origin selects a host and is required only with more than one host.
// Who is an email or existing username; an unknown email creates a user.
// Both print a secret, one-use sign-in link that expires in 30 minutes.
// Reset unlinks Google and revokes passkeys, sessions, codes and links on all hosts.
//
// # Usernames and email
//
// Usernames default to email; older accounts get one at their next sign-in.
// A username contains @ only when it equals the user's own email, ignoring case.
// Operator commands accept usernames even for users without an email.
// Email changes need a mail queue and a link sent to the new address.
// The change revokes codes and links for both addresses and notifies the old one.
//
// # Styling and browser support
//
// Auth ships no CSS; the site styles these hooks:
//   - .auth[data-auth]: page roots (signin, code, link, email, account, continue or forbidden).
//   - .auth-field: field labels, including the delete confirmation.
//   - .auth-primary, .auth-secondary and .auth-danger: action buttons.
//   - .auth-note: supporting text; .auth-problem: error messages.
//   - .auth-passkeys: the passkey list.
//
// Auth serves passkey.js and loads it on the pages that need it.
// The site's CSP must allow script-src 'self' and connect-src 'self'.
// Passkeys need Chrome 129, Safari 18.4 or Firefox 119 or later.
// Older browsers can use email codes in Invite and Open modes.
//
// # Overriding a page
//
// Set Host.Views.Signin to a func(SigninView) templ.Component to replace sign-in.
// Unset Views fields keep the defaults; start from the body in views.templ.
// Keep the model's feature and safety guards, such as Codes and KeepLast.
// Keep form methods, actions, field names and hidden next values.
// Keep data-auth-passkey="login" or "register" buttons and their data-next.
// Keep data-auth-link-token and data-auth-problem hooks.
// Keep autocomplete="username webauthn" on the email input for passkey autofill.
// Google sign-in and linking forms need hx-boost="false" for full-page navigation.
// Return only the body, with no data-auth root or script tag of its own.
// Auth owns the root, script loading and Host.Layout call.
package auth

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"
	"slices"
	"strings"
	"time"

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
	// Invite adds emailed codes and links for existing users.
	Invite
	// Open creates users on code or link sign-in and allows account deletion.
	Open
)

func (m Mode) codes() bool { return m != Operator }

// Config describes shared accounts and their host policies.
type Config struct {
	// Name is the site name shown by passkey managers and in mail.
	Name string
	// Key is the HMAC key for codes and links, at least 32 bytes.
	Key []byte
	// Google enables Google sign-in when non-nil.
	Google *Google
	// Hosts must share a registrable domain or use localhost names.
	// A single loopback IP is also allowed; at least one host is required.
	Hosts []Host
}

// User is an account shared by all configured hosts.
type User struct {
	ID       pg.ID
	Email    string // Mailbox-verified or operator-supplied address; "" if none.
	Username string
	Created  time.Time
}

// Auth serves sign-in, the account page and session lookups.
type Auth struct {
	db      *gorm.DB
	mail    *mail.Queue
	name    string
	key     []byte
	hosts   map[string]*host
	rp      string
	wa      *webauthn.WebAuthn
	assets  *static.Server
	script  string
	library string
	google  *provider
}

const (
	prefix       = "/auth/"
	signinPath   = prefix + "signin"
	accountPath  = prefix + "account"
	confirmPath  = prefix + "email"
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

// Anyone adds no route-level restriction to Require.
func Anyone(*http.Request, User) bool { return true }

// New checks c and prepares handlers; it does not migrate the database.
// A non-nil c.Google needs both ClientID and ClientSecret supplied by the site.
// A nil q is allowed only when every host is Operator and disables mail notices.
// The caller owns q and must close it on shutdown.
func New(db *gorm.DB, q *mail.Queue, c Config) (*Auth, error) {
	hosts, rp, err := c.check()
	if err != nil {
		return nil, err
	}
	origins := make([]string, 0, len(hosts))
	for _, h := range hosts {
		if h.layout == nil {
			return nil, fmt.Errorf("auth: host %s needs a Layout", h.origin)
		}
		if h.mode.codes() && q == nil {
			return nil, fmt.Errorf("auth: host %s: Invite and Open modes need a mail queue", h.origin)
		}
		origins = append(origins, h.origin)
	}
	slices.Sort(origins)
	wa, err := webauthn.New(&webauthn.Config{
		RPID:                  rp,
		RPDisplayName:         c.Name,
		RPOrigins:             origins,
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
	var google *provider
	if c.Google != nil {
		google = newProvider(googleKind, googleIssuer, c.Google.ClientID, c.Google.ClientSecret)
	}
	return &Auth{
		db:      db,
		mail:    q,
		name:    c.Name,
		key:     bytes.Clone(c.Key),
		hosts:   hosts,
		rp:      rp,
		wa:      wa,
		assets:  assets,
		script:  script,
		library: library(),
		google:  google,
	}, nil
}

func (c Config) check() (map[string]*host, string, error) {
	if c.Name == "" {
		return nil, "", errors.New("auth: Config.Name is required")
	}
	if len(c.Key) < 32 {
		return nil, "", errors.New("auth: Config.Key must be at least 32 bytes")
	}
	if c.Google != nil && (c.Google.ClientID == "" || c.Google.ClientSecret == "") {
		return nil, "", errors.New("auth: Config.Google needs a ClientID and ClientSecret")
	}
	return parseHosts(c.Hosts)
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
	handle := func(pattern string, serve func(http.ResponseWriter, *http.Request, *host)) {
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			h, ok := a.host(r)
			if !ok {
				misdirected(w)
				return
			}
			serve(w, r, h)
		})
	}
	mux.Handle("GET "+assetsPath, a.assets)
	handle("GET "+signinPath, a.signin)
	handle("POST /auth/code", a.sendCode)
	handle("POST /auth/code/verify", a.verifyCode)
	handle("POST /auth/passkey/begin", a.loginBegin)
	handle("POST /auth/passkey/finish", a.loginFinish)
	handle("GET /auth/link", a.showLink)
	handle("POST /auth/link", a.redeemLink)
	handle("POST /auth/signout", a.signout)
	handle("GET "+accountPath, a.account)
	handle("POST /auth/account/passkeys/begin", a.addBegin)
	handle("POST /auth/account/passkeys/finish", a.addFinish)
	handle("POST /auth/account/passkeys/{id}/rename", a.rename)
	handle("POST /auth/account/passkeys/{id}/remove", a.remove)
	handle("POST /auth/account/signout-others", a.signoutOthers)
	handle("POST /auth/account/username", a.setUsername)
	handle("POST /auth/account/delete", a.deleteAccount)
	if a.mail != nil {
		handle("POST /auth/account/email", a.changeEmail)
		handle("GET "+confirmPath, a.showConfirm)
		handle("POST "+confirmPath, a.confirmEmail)
	}
	if a.google != nil {
		handle("POST /auth/google", a.googleSignin)
		handle("GET "+googleCallback, a.googleCallback)
		handle("POST /auth/account/google", a.googleLink)
		handle("POST /auth/account/google/remove", a.googleUnlink)
	}
}

type sessionKey struct{}

type cached struct {
	host    *host
	session session
}

// Require serves handler to signed-in users that allow accepts; others sign in or get 403.
// Pass Anyone to add no route-level restriction; Host.Allow still applies.
func (a *Auth) Require(allow func(*http.Request, User) bool, handler http.Handler) http.Handler {
	if allow == nil {
		panic("auth: Require needs an allow check; use auth.Anyone for any signed-in user")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, ok := a.host(r)
		if !ok {
			misdirected(w)
			return
		}
		s, ok := a.session(r, h)
		if !ok {
			http.Redirect(w, r, signinPath+"?"+url.Values{"next": {r.URL.RequestURI()}}.Encode(), http.StatusSeeOther)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), sessionKey{}, cached{h, s}))
		if !h.admits(s.user) || !allow(r, s.user) {
			a.forbid(w, r, h, s.user)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// User returns the signed-in user, reusing Require's lookup when present.
func (a *Auth) User(r *http.Request) (User, bool) {
	h, ok := a.host(r)
	if !ok {
		return User{}, false
	}
	s, ok := a.session(r, h)
	if !ok || !h.admits(s.user) {
		return User{}, false
	}
	return s.user, true
}
