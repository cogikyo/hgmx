package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/a-h/templ"

	"github.com/cogikyo/hgmx/render"
)

func (a *Auth) page(w http.ResponseWriter, r *http.Request, status int, title string, body templ.Component) {
	w.Header().Set("Cache-Control", "no-store")
	err := render.HTML(w, r, status, a.c.Layout(r, title, body))
	if err == nil || r.Context().Err() != nil {
		return
	}
	panic(err)
}

func (a *Auth) signin(w http.ResponseWriter, r *http.Request) {
	a.page(w, r, http.StatusOK, "Sign in", signinPage(a.signinView(r, next(r.URL.Query().Get("next")), "", "")))
}

func (a *Auth) toSignin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, signinPath+"?"+url.Values{"next": {accountPath}}.Encode(), http.StatusSeeOther)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

func problem(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"problem": msg})
}

func next(raw string) string {
	p, _, _ := strings.Cut(raw, "?")
	if _, err := url.Parse(raw); err != nil || !strings.HasPrefix(raw, "/") || strings.Contains(raw, `\`) || strings.Contains(p, "//") {
		return accountPath
	}
	return raw
}
