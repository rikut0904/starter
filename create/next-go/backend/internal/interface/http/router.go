package http

import (
	"github.com/labstack/echo/v5"
	"github.com/rikut0904/starter/create/next-go/backend/internal/usecase"
)

func NewRouter(health usecase.Health) *echo.Echo {
	e := echo.New()
	e.GET("/health", NewHandler(health))
	return e
}
