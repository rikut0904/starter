package http

import (
	"github.com/labstack/echo/v5"
	"github.com/rikut0904/starter/create/go/internal/usecase"
)

func NewHandler(health usecase.Health) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.JSON(200, health.Check())
	}
}
