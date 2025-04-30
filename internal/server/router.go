package server

import (
	// "net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewRouter creates and configures a new Echo instance.
func NewRouter() *echo.Echo {
	e := echo.New()

	// --- Middleware ---
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- Static Files ---
	e.Static("/static", "static")

	// --- Application Routes ---
	e.GET("/", Index)

	return e
}
