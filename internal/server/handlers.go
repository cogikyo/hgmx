package server

import (

	// "net/http"

	"github.com/labstack/echo/v4"
	"github.com/nosvagor/hgmx/views/builder"
	"github.com/nosvagor/hgmx/views/index"
)

// Index handler for the root path using Echo context.
func Index(c echo.Context) error {
	cmp := index.Index(index.Page{
		Title:   "HGMX Builder",
		Content: builder.Palette(),
	})

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
