package server

import (
	"context"
	"fmt"
	"log/slog"

	"baneks.com/internal/api/middlewares"
	customerrors "baneks.com/internal/custom_errors"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func InitializeServer(
	ctx context.Context,
	logger *slog.Logger,
	guardApiKey string,
) *echo.Echo {
	server := echo.New()
	initMiddlewares(ctx, server, logger, guardApiKey)

	server.Logger = slog.Default()
	validator := CreateValidator()
	server.Validator = validator

	return server
}

func errorAttrs(e error) []slog.Attr {
	attrs := []slog.Attr{}

	switch v := e.(type) {
	case *customerrors.AppHttpError:
		var msg string
		errBytes, err := v.MarshalJSON()
		if err != nil {
			// если сериализация упала — выводим fallback
			msg = fmt.Sprintf("Failed to marshal AppHttpError, code=%d: %v", v.Code, err)
		} else {
			msg = string(errBytes)
		}

		attrs = append(attrs, slog.String("err_message", msg))
		attrs = append(attrs, slog.Int("err_code", v.Code))

		if v.Internal != nil {
			attrs = append(attrs, slog.String("err_internal", v.Internal.Error()))
		}

	default:
		attrs = append(attrs, slog.String("err_message", e.Error()))
	}

	return attrs
}

func initMiddlewares(
	ctx context.Context,
	e *echo.Echo,
	logger *slog.Logger,
	guardApiKey string,
) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogMethod: true,
		LogStatus: true,

		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
			}

			var level slog.Level
			var msg string

			switch {
			case v.Status >= 500:
				level = slog.LevelError
				msg = "REQUEST_SERVER_ERROR"
			case v.Status >= 400:
				level = slog.LevelInfo
				msg = "REQUEST_CLIENT_ERROR"
			default:
				level = slog.LevelInfo
				msg = "REQUEST_OK"
			}

			if v.Error != nil {
				attrs = append(attrs, errorAttrs(v.Error)...)
			}

			logger.LogAttrs(ctx, level, msg, attrs...)
			return nil
		},
	}))

	e.Use(middleware.Recover())

	// global middlewares init
	guard := middlewares.NewGuardWithSecret(guardApiKey)
	e.Use(guard.GuardWithSecretMiddleware)
}
