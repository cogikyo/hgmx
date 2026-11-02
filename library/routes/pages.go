package routes

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/cogikyo/hgmx/render"
	"github.com/cogikyo/hgmx/server"

	"github.com/cogikyo/hgmx/library/store"
	"github.com/cogikyo/hgmx/library/views"
)

type pages struct {
	store *store.Store
	views *views.Site
}

func (p *pages) page(status int, title string, content templ.Component) http.Handler {
	return server.PageVersion(p.views.Version, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := render.HTML(w, r, status, p.views.Page(r.URL.Path, title, content))
		if err == nil || r.Context().Err() != nil {
			return
		}
		panic(err)
	}))
}
