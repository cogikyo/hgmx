package server

import (

	// "net/http"

	"log"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/nosvagor/hgmx/internal/palette"
	"github.com/nosvagor/hgmx/views/builder"
	"github.com/nosvagor/hgmx/views/index"
)

// Index handler for the root path using Echo context.
func Index(c echo.Context) error {
	cmp := index.Index(index.Page{
		Title:   "HGMX Builder",
		Content: Palette(),
	})

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func Palette() templ.Component {
	hex := "#222536"
	oklch, _ := palette.HexToOklch(hex)
	p := palette.Generate(oklch)
	viewModel := p.ToView()
	log.Println(viewModel)
	return builder.Palette(viewModel)
}
