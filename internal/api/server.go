package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeServer() *echo.Echo {
	server := echo.New()
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	server.HideBanner = true
	server.HTTPErrorHandler = CreateErrorHandler()
	server.Validator = CreateValidator()
	return server
}
