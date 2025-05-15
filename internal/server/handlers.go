package server

import (
	"github.com/labstack/echo/v4"
	views "github.com/nosvagor/hgmx/app/views"
	colors "github.com/nosvagor/hgmx/app/views/pages/colors"
	"github.com/nosvagor/hgmx/internal/palette"
)

// Index handler for the root path using Echo context.
func Index(c echo.Context) error {
	cmp := views.FullPage(views.Page{
		Title: "Index",
	})

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func Palette(c echo.Context) error {
	hex := c.QueryParam("hex")
	if hex == "" {
		hex = "#222536"
	}
	// hex = "#ffffff"
	p := palette.Generate(hex)
	viewModel := p.ToView()

	cmp := views.FullPage(views.Page{
		Title:   "Palette View",
		Content: colors.Palette(viewModel, hex),
	})

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
