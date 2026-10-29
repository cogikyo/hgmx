// Package views wires the scaffold's templates to hgmx asset handlers.
package views

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/cogikyo/hgmx/render"
	"github.com/cogikyo/hgmx/scripts"
	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/static"
)

//go:embed static
var files embed.FS

const assets = "/static/"

type site struct {
	static  *static.Server
	scripts *scripts.Handler
	version string
}

func Mount(mux *http.ServeMux, version, dir string) (fallback http.Handler, err error) {
	s := &site{scripts: scripts.New(scripts.Navigation), version: version}
	if dir != "" {
		s.static = static.Dev(assets, os.DirFS(dir))
	} else {
		sub, err := fs.Sub(files, "static")
		if err != nil {
			return nil, err
		}
		if s.static, err = static.New(assets, sub); err != nil {
			return nil, err
		}
	}

	mux.Handle("GET "+assets, s.static)
	mux.Handle("GET "+scripts.Path, s.scripts)
	mux.Handle("GET /{$}", s.page("/", "Home", home()))
	mux.Handle("GET /about", s.page("/about", "About", about()))
	mux.Handle("GET /", server.PageVersion(version, s.failure(http.StatusNotFound)))
	return s.failure(http.StatusInternalServerError), nil
}

func (s *site) page(path, title string, content templ.Component) http.Handler {
	return server.PageVersion(s.version, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, http.StatusOK, full(s, path, title, content))
	}))
}

func (s *site) failure(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, status, full(s, "", http.StatusText(status), failure(status)))
	})
}

func (s *site) render(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	err := render.HTML(w, r, status, c)
	if err == nil {
		return
	}
	slog.ErrorContext(r.Context(), "render", "path", r.URL.Path, "err", err)
	if status == http.StatusInternalServerError {
		http.Error(w, http.StatusText(status), status)
		return
	}
	s.failure(http.StatusInternalServerError).ServeHTTP(w, r)
}
