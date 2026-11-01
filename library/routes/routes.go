package routes

import (
	"net/http"

	"github.com/cogikyo/hgmx/scripts"

	"github.com/cogikyo/hgmx/library/store"
	"github.com/cogikyo/hgmx/library/views"
	"github.com/cogikyo/hgmx/library/views/pages/about"
	"github.com/cogikyo/hgmx/library/views/pages/home"
	"github.com/cogikyo/hgmx/library/views/pages/notfound"
)

func New(st *store.Store, v *views.Site) *http.ServeMux {
	p := &pages{store: st, views: v}

	mux := http.NewServeMux()
	mux.Handle("GET "+views.Assets, v.Static)
	mux.Handle("GET "+scripts.Path, v.Scripts)
	mux.Handle("GET /favicon.ico", v.Static.File("favicon/favicon.ico"))
	mux.Handle("GET /{$}", p.page(http.StatusOK, "Home", home.Page()))
	mux.Handle("GET /about", p.page(http.StatusOK, "About", about.Page()))
	mux.Handle("GET /", p.page(http.StatusNotFound, "Not Found", notfound.Page()))
	return mux
}
