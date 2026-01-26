package server

import (
	"context"
	"log/slog"
	"net/http"

	customerrors "baneks.com/internal/custom_errors"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func InitializeServer(
	ctx context.Context,
	logger *slog.Logger,
) *echo.Echo {
	server := echo.New()
	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogMethod: true,
		LogStatus: true,

		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				logger.LogAttrs(ctx,
					slog.LevelError,
					"REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", getError(v.Error)),
				)
			} else {
				logger.LogAttrs(
					ctx,
					slog.LevelInfo,
					"REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			}
			return nil
		},
	}))
	server.Use(middleware.Recover())
	server.Logger = slog.Default()
	validator := CreateValidator()
	server.Validator = validator

	return server
}

func getError(e error) string {
	switch v := e.(type) {
	case *customerrors.AppHttpError:
		if v.Internal != nil {
			return v.Internal.Error()
		} else if v.Message != "" {
			return v.MessageString()
		}
		return http.StatusText(v.Code)
	default:
		return e.Error()
	}
}
