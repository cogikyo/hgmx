// Package views wires the scaffold's templates to hgmx asset handlers.
package views

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"

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
	head    template.HTML
	version string
}

type Config struct {
	Version string
	Static  string        // Development directory; empty uses embedded assets.
	Head    template.HTML // Trusted HTML appended to the document head.
}

// New builds the application mux; wrap it with [server.Site.Handler].
func New(c Config) (*http.ServeMux, error) {
	s := &site{scripts: scripts.New(scripts.Navigation), head: c.Head, version: c.Version}
	if c.Static != "" {
		s.static = static.Dev(assets, os.DirFS(c.Static))
	} else {
		sub, err := fs.Sub(files, "static")
		if err != nil {
			return nil, err
		}
		if s.static, err = static.New(assets, sub); err != nil {
			return nil, err
		}
	}

	mux := http.NewServeMux()
	mux.Handle("GET "+assets, s.static)
	mux.Handle("GET "+scripts.Path, s.scripts)
	mux.Handle("GET /{$}", s.page(http.StatusOK, "Home", home()))
	mux.Handle("GET /about", s.page(http.StatusOK, "About", about()))
	mux.Handle("GET /", s.page(http.StatusNotFound, "Not Found", notFound()))
	return mux, nil
}

func (s *site) page(status int, title string, content templ.Component) http.Handler {
	return server.PageVersion(s.version, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := render.HTML(w, r, status, full(s, r.URL.Path, title, content))
		if err == nil || r.Context().Err() != nil {
			return
		}
		panic(err)
	}))
}

const internalError = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Internal Server Error</title>
</head>
<body>
<main id="main">
<h1>Internal Server Error</h1>
<p>Something went wrong on our side. <a href="/">Go home</a>.</p>
</main>
</body>
</html>
`

func InternalError(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Cache-Control", "no-store")
	h.Set("Content-Length", strconv.Itoa(len(internalError)))
	w.WriteHeader(http.StatusInternalServerError)
	if r.Method != http.MethodHead {
		io.WriteString(w, internalError)
	}
}
