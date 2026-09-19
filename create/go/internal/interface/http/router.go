package http

import (
	"github.com/labstack/echo/v5"
	"github.com/rikut0904/starter/create/go/internal/usecase"
)

// NewRouter wires HTTP routes in one place so handlers remain independent of transport setup.
func NewRouter(health usecase.Health) *echo.Echo {
	e := echo.New()
	e.GET("/health", NewHandler(health))
	return e
}
