package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeServer() *echo.Echo {
	server := echo.New()
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "time: ${time_custom}: method=${method}, uri=${uri}, status=${status}, error=${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00",
	}))
	server.HideBanner = true
	server.HTTPErrorHandler = CreateErrorHandler()
	validator := CreateValidator()
	server.Validator = validator

	return server
}
