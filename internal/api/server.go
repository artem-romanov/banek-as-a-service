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
				attrs := []slog.Attr{
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				}

				attrs = append(attrs, errorAttrs(v.Error)...)

				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"REQUEST_ERROR",
					attrs...,
				)
			} else {
				logger.LogAttrs(
					ctx,
					slog.LevelInfo,
					"REQUEST",
					slog.String("method", v.Method),
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

func errorAttrs(e error) []slog.Attr {
	attrs := []slog.Attr{}

	switch v := e.(type) {
	case *customerrors.AppHttpError:
		if msg := v.MessageString(); msg != "" {
			attrs = append(attrs, slog.String("err_message", msg))
		} else {
			attrs = append(attrs, slog.String("err_message", http.StatusText(v.Code)))
		}

		if v.Internal != nil {
			attrs = append(attrs, slog.String("err_internal", v.Internal.Error()))
		}

	default:
		attrs = append(attrs, slog.String("err_message", e.Error()))
	}

	return attrs
}
