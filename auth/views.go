package auth

import (
	"net/http"
	"strconv"
	"time"
)

type signinView struct {
	Script   string
	Next     string
	Codes    bool
	Email    string
	SignedIn string
	Problem  string
}

type codeView struct {
	Email   string
	Next    string
	Open    bool
	Problem string
}

type linkView struct {
	Script  string
	Next    string
	Problem string
}

type passkeyRow struct {
	ID      string
	Name    string
	Created time.Time
	Used    *time.Time
}

type accountView struct {
	Script   string
	Email    string
	Next     string
	Passkeys []passkeyRow
	Others   int
	Recent   bool
	KeepLast bool
	Delete   bool
	Problem  string
}

func (a *Auth) signinView(r *http.Request, next, email, problem string) signinView {
	v := signinView{Script: a.script, Next: next, Codes: a.c.Mode.codes(), Email: email, Problem: problem}
	if s, ok := a.session(r); ok {
		v.SignedIn = s.user.Email
	}
	return v
}

func date(t time.Time) string { return t.UTC().Format("2 Jan 2006") }

func used(t *time.Time) string {
	if t == nil {
		return "Not used yet."
	}
	return "Last used " + date(*t) + "."
}

func devices(n int) string {
	if n == 1 {
		return "1 other device"
	}
	return strconv.Itoa(n) + " other devices"
}
