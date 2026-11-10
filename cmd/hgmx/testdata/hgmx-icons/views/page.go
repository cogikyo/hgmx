package views

import (
	catalog "github.com/cogikyo/hgmx/icons"

	"site/views/icons"
)

var (
	_ = catalog.New()
	_ = catalog.Path
	_ = []any{catalog.UserCircle, catalog.TriangleAlert, catalog.Plus, catalog.Hover}
	_ = icons.UserCircle
)
