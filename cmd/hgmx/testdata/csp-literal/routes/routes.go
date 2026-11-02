package routes

import (
	"net/http"

	"github.com/cogikyo/hgmx/server"
)

const csp = "Content-Security-Policy"

type policy struct{}

func (policy) Header() string { return "sandbox" }

func Files(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
	h := w.Header()
	h.Add("content-security-policy", "default-src 'none'")
	h[csp] = []string{"default-src 'none'"}
	h.Set(csp, policy{}.Header())
	h.Set(csp, server.CSP{"sandbox": nil}.Header())
	h[csp] = []string{server.CSP{"sandbox": nil}.Header()}
	_ = http.Header{csp: {"sandbox"}}
	_ = http.Header{csp: {server.CSP{"sandbox": nil}.Header()}, "Cache-Control": {"no-store"}}
	h.Set("Content-Security-Policy-Report-Only", "default-src 'none'")
	h.Set("Cache-Control", "no-store")
}
