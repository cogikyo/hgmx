// Package views wires the scaffold's templates to hgmx asset handlers.
package views

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/cogikyo/hgmx/scripts"
	"github.com/cogikyo/hgmx/static"
	"github.com/cogikyo/hgmx/webapp"
)

//go:embed static
var files embed.FS

//go:embed icon.png
var icon []byte

const assets = "/static/"

// Site groups the scaffold's static, script and app handlers.
type Site struct {
	// Static serves site assets, not vendored scripts or app icons.
	Static  *static.Server
	Scripts *scripts.Handler
	// App is initialized from the embedded icon.png by [New].
	App *webapp.Handler
}

// New serves embedded assets, or views/static from disk in dev mode.
// Scripts and the app icon still come from embedded files.
func New(dev bool) (*Site, error) {
	var site Site
	var err error
	if dev {
		site.Static = static.Dev(assets, os.DirFS("views/static"))
	} else {
		sub, err := fs.Sub(files, "static")
		if err != nil {
			return nil, err
		}
		if site.Static, err = static.New(assets, sub); err != nil {
			return nil, err
		}
	}
	if site.Scripts, err = scripts.New(scripts.Preload, scripts.Hyperscript); err != nil {
		return nil, err
	}
	site.App, err = webapp.New(webapp.App{
		Name:       "hgmx",
		ShortName:  "hgmx",
		Theme:      "#010a05",
		Background: "#010a05",
		Icon:       icon,
	})
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// Mount claims GET /; register caller routes with a method.
func (s *Site) Mount(mux *http.ServeMux) {
	mux.Handle("GET "+assets, s.Static)
	mux.Handle("GET "+scripts.Path, s.Scripts)
	mux.Handle("GET "+webapp.Path, s.App)
	mux.Handle("GET /{$}", templ.Handler(Full(s, "/", "Home", Home())))
	mux.Handle("GET /about", templ.Handler(Full(s, "/about", "About", About())))
	mux.Handle("GET /", templ.Handler(Full(s, "", "Not Found", NotFound()), templ.WithStatus(http.StatusNotFound)))
}
