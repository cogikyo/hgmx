package routes

import (
	"html/template"
	"net/http"

	"github.com/cogikyo/hgmx/scripts"

	"github.com/cogikyo/hgmx/library/store"
	"github.com/cogikyo/hgmx/library/views"
	"github.com/cogikyo/hgmx/library/views/pages/about"
	"github.com/cogikyo/hgmx/library/views/pages/home"
	"github.com/cogikyo/hgmx/library/views/pages/notfound"
)

type Config struct {
	Store   *store.Store
	Version string
	Static  string
	Head    template.HTML
}

func New(c Config) (*http.ServeMux, error) {
	v, err := views.New(views.Config{Version: c.Version, Static: c.Static, Head: c.Head})
	if err != nil {
		return nil, err
	}
	p := &pages{store: c.Store, views: v, version: c.Version}

	mux := http.NewServeMux()
	mux.Handle("GET "+views.Assets, v.Static)
	mux.Handle("GET "+scripts.Path, v.Scripts)
	mux.Handle("GET /{$}", p.page(http.StatusOK, "Home", home.Page()))
	mux.Handle("GET /about", p.page(http.StatusOK, "About", about.Page()))
	mux.Handle("GET /", p.page(http.StatusNotFound, "Not Found", notfound.Page()))
	return mux, nil
}
