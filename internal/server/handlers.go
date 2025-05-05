package server

import (

	// "net/http"

	"github.com/labstack/echo/v4"
	"github.com/nosvagor/hgmx/internal/palette"
	"github.com/nosvagor/hgmx/views"
	"github.com/nosvagor/hgmx/views/builder"
	"github.com/nosvagor/hgmx/views/testing"
)

// Index handler for the root path using Echo context.
func Index(c echo.Context) error {
	return Palette(c)
}

func Palette(c echo.Context) error {
	hex := "#222536"
	oklch, _ := palette.HexToOklch(hex)
	p := palette.Generate(oklch)
	viewModel := p.ToView()

	cmp := views.FullPage(views.Page{
		Title:   "Palette View",
		Content: builder.Palette(viewModel),
	})

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func Testing(c echo.Context) error {
	cmp := views.FullPage(views.Page{
		Title:   "Testing",
		Content: testing.Main(),
	})

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
