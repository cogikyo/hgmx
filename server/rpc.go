package server

import (
	"net/http"
	"strings"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
)

type rpc struct {
	services map[string]bool
	errors   *connecthttp.ErrorWriter
}

func newRPC(s *connect.Server) *rpc {
	p := &rpc{services: map[string]bool{}, errors: connecthttp.NewErrorWriter()}
	for spec := range s.Specs() {
		if i := strings.LastIndexByte(spec.Procedure, '/'); i > 0 {
			p.services[spec.Procedure[:i+1]] = true
		}
	}
	return p
}

func (p *rpc) known(path string) bool {
	i := strings.LastIndexByte(path, '/')
	return i > 0 && p.services[path[:i+1]]
}

func (p *rpc) shaped(r *http.Request) bool {
	return p.known(r.URL.Path) ||
		r.Header.Get("Connect-Protocol-Version") != "" ||
		strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc")
}

func (p *rpc) route(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !p.known(r.URL.Path) && p.shaped(r) {
			w.Header().Set("Cache-Control", "no-store")
			p.errors.Write(w, r, connect.NewError(connect.CodeUnimplemented, "no such procedure"))
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (p *rpc) fallback(page http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !p.shaped(r) {
			page.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		p.errors.Write(w, r, connect.NewError(connect.CodeInternal, "Something went wrong. Try again."))
	})
}
