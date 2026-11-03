package auth

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
)

// Views replaces page bodies per host; nil fields use [DefaultViews].
type Views struct {
	Signin    func(SigninView) templ.Component
	Code      func(CodeView) templ.Component
	Link      func(LinkView) templ.Component
	Email     func(LinkView) templ.Component
	Account   func(AccountView) templ.Component
	Continue  func(ContinueView) templ.Component
	Forbidden func(ForbiddenView) templ.Component
}

// DefaultViews returns page bodies without roots, scripts or layouts.
func DefaultViews() Views {
	return Views{
		Signin:    signinBody,
		Code:      codeBody,
		Link:      linkBody,
		Email:     emailBody,
		Account:   accountBody,
		Continue:  continueBody,
		Forbidden: forbiddenBody,
	}
}

func (v Views) withDefaults() Views {
	d := DefaultViews()
	return Views{
		Signin:    or(v.Signin, d.Signin),
		Code:      or(v.Code, d.Code),
		Link:      or(v.Link, d.Link),
		Email:     or(v.Email, d.Email),
		Account:   or(v.Account, d.Account),
		Continue:  or(v.Continue, d.Continue),
		Forbidden: or(v.Forbidden, d.Forbidden),
	}
}

func or[V any](view, fallback func(V) templ.Component) func(V) templ.Component {
	if view == nil {
		return fallback
	}
	return view
}

type SigninView struct {
	Next     string
	Codes    bool
	Email    string
	SignedIn string // Display name, not a user ID.
	Google   bool
	Problem  string
}

type CodeView struct {
	Email   string
	Next    string
	Open    bool
	Problem string
}

// LinkView serves both sign-in links and email confirmation.
type LinkView struct {
	Next    string
	Problem string
}

type PasskeyRow struct {
	ID      string
	Name    string
	Created time.Time
	Used    *time.Time
}

type SessionRow struct {
	ID      string
	Device  string
	Seen    time.Time
	Current bool
}

type AccountView struct {
	Name       string
	Username   string
	Email      string
	Mail       bool
	Next       string
	Passkeys   []PasskeyRow // Usable with the current RP ID.
	Older      []PasskeyRow // Remove-only keys unusable with this RP ID.
	Sessions   []SessionRow // Browser and app sessions here, current included.
	Recent     bool         // Signed in within the last 10 minutes.
	KeepLast   bool         // The sole sign-in method is a passkey.
	Google     bool
	Linked     bool
	KeepGoogle bool // Google is the sole sign-in method.
	Delete     bool
	Problem    string
	Notice     string
}

type ContinueView struct {
	URL string
}

type ForbiddenView struct {
	Name string // Display name of the refused user.
}

func (a *Auth) signinPage(w http.ResponseWriter, r *http.Request, h *host, status int, next, email, problem string) {
	v := SigninView{Next: next, Codes: h.mode.codes(), Email: email, Google: a.google != nil, Problem: problem}
	if s, ok := a.session(r, h); ok && h.admits(s.user) {
		v.SignedIn = s.user.name()
	}
	a.page(w, r, h, status, "Sign in", root("signin", a.script, h.views.Signin(v)))
}

func (a *Auth) codePage(w http.ResponseWriter, r *http.Request, h *host, status int, v CodeView) {
	a.page(w, r, h, status, "Enter your code", root("code", "", h.views.Code(v)))
}

func (a *Auth) linkPage(w http.ResponseWriter, r *http.Request, h *host, status int, next, problem string) {
	a.page(w, r, h, status, "Sign in", root("link", a.script, h.views.Link(LinkView{Next: next, Problem: problem})))
}

func (a *Auth) emailPage(w http.ResponseWriter, r *http.Request, h *host, status int, problem string) {
	a.page(w, r, h, status, "Confirm email", root("email", a.script, h.views.Email(LinkView{Problem: problem})))
}

func (a *Auth) continuePage(w http.ResponseWriter, r *http.Request, h *host, to string) {
	a.page(w, r, h, http.StatusOK, "Continue to Google", continueRoot(to, h.views.Continue(ContinueView{URL: to})))
}

func (a *Auth) forbid(w http.ResponseWriter, r *http.Request, h *host, u User) {
	a.page(w, r, h, http.StatusForbidden, "Not allowed", root("forbidden", "", h.views.Forbidden(ForbiddenView{Name: u.name()})))
}

func date(t time.Time) string { return t.UTC().Format("2 Jan 2006") }

func used(t *time.Time) string {
	if t == nil {
		return "Not used yet."
	}
	return "Last used " + date(*t) + "."
}
