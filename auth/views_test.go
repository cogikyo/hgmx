package auth

import (
	"context"
	"html"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func brand(page, problem string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<h1 class="brand">`+page+`</h1><p role="alert" data-auth-problem>`+html.EscapeString(problem)+`</p>`)
		return err
	})
}

func TestCustomViews(t *testing.T) {
	a := setup(t)
	c := config()
	c.Hosts[0].Layout = func(_ *http.Request, title string, body templ.Component) templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			io.WriteString(w, "<!doctype html><title>"+title+"</title><main>")
			if err := body.Render(ctx, w); err != nil {
				return err
			}
			_, err := io.WriteString(w, "</main>")
			return err
		})
	}
	c.Hosts[0].Views = Views{
		Signin:  func(v SigninView) templ.Component { return brand("signin "+v.Next, v.Problem) },
		Account: func(v AccountView) templ.Component { return brand("account "+v.Username, v.Problem) },
	}
	a, err := New(a.db, nil, c)
	if err != nil {
		t.Fatal(err)
	}
	at(a, apex).mode = Invite
	mux := http.NewServeMux()
	a.Mount(mux)

	u, err := createUser(a.db, "brand@example.com")
	if err != nil {
		t.Fatal(err)
	}
	token, err := startSession(a.db, httptest.NewRequest("POST", "/", nil), at(a, apex), u.ID)
	if err != nil {
		t.Fatal(err)
	}
	serve := func(method, target string, form url.Values) (int, string) {
		r := httptest.NewRequest(method, target, strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.AddCookie(&http.Cookie{Name: sessionCookie, Value: token})
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w.Code, w.Body.String()
	}

	for _, c := range []struct {
		name, method, target string
		form                 url.Values
		status               int
		root, want           string
		script, document     bool
	}{
		{"custom signin", "GET", apex + "/auth/signin?next=/shop", nil, 200, "signin", `<h1 class="brand">signin /shop</h1>`, true, true},
		{"custom signin on a validation failure", "POST", apex + "/auth/code", url.Values{"email": {"nope"}, "next": {"//evil.example"}}, 422, "signin", `signin /auth/account</h1><p role="alert" data-auth-problem>Enter a valid email address.`, true, true},
		{"default code page", "POST", apex + "/auth/code/verify", url.Values{"email": {"a@example.com"}, "code": {"x"}}, 422, "code", `<h1>Check your email</h1>`, false, true},
		{"custom account", "GET", apex + "/auth/account", nil, 200, "account", `<h1 class="brand">account brand@example.com</h1>`, true, true},
		{"custom account on a validation failure", "POST", apex + "/auth/account/username", url.Values{"username": {"a b"}}, 422, "account", `data-auth-problem>Usernames need`, true, true},
		{"defaults on the other host", "GET", admin + "/auth/signin", nil, 200, "signin", `<h1>Sign in</h1>`, true, false},
	} {
		status, body := serve(c.method, c.target, c.form)
		if status != c.status || !strings.Contains(body, c.want) {
			t.Errorf("%s: status %d, body %s; want %d containing %q", c.name, status, body, c.status, c.want)
			continue
		}
		if n := strings.Count(body, "data-auth="); n != 1 || !strings.Contains(body, `<div class="auth" data-auth="`+c.root+`">`) {
			t.Errorf("%s: %d data-auth roots, want one %q root: %s", c.name, n, c.root, body)
		}
		if got := strings.Count(body, `<script type="module" src="`+a.script+`">`); got != map[bool]int{false: 0, true: 1}[c.script] {
			t.Errorf("%s: passkey.js loaded %d times, want script %v", c.name, got, c.script)
		}
		if c.document && !strings.HasPrefix(body, "<!doctype html>") {
			t.Errorf("%s: response is not the site's Layout document: %s", c.name, body)
		}
	}
}
