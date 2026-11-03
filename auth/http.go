package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/a-h/templ"

	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/render"
)

func (a *Auth) page(w http.ResponseWriter, r *http.Request, h *host, status int, title string, body templ.Component) {
	w.Header().Set("Cache-Control", "no-store")
	err := render.HTML(w, r, status, h.layout(r, title, body))
	if err == nil || r.Context().Err() != nil {
		return
	}
	panic(err)
}

func (a *Auth) signin(w http.ResponseWriter, r *http.Request, h *host) {
	a.signinPage(w, r, h, http.StatusOK, next(r.URL.Query().Get("next")), "", "")
}

func (a *Auth) toSignin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, signinPath+"?"+url.Values{"next": {accountPath}}.Encode(), http.StatusSeeOther)
}

func (a *Auth) signedIn(w http.ResponseWriter, r *http.Request, h *host) (session, bool) {
	s, ok := a.session(r, h)
	switch {
	case !ok:
		a.toSignin(w, r)
	case !h.admits(s.user):
		a.forbid(w, r, h, s.user)
	default:
		return s, true
	}
	return s, false
}

func (a *Auth) refused(w http.ResponseWriter, r *http.Request, h *host, s session, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, errStale), errors.Is(err, pg.ErrNotFound):
		a.toSignin(w, r)
	case errors.Is(err, errDenied):
		a.forbid(w, r, h, s.user)
	default:
		panic(err)
	}
	return true
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
