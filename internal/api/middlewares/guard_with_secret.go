package middlewares

import (
	"errors"
	"net/http"

	customerrors "baneks.com/internal/custom_errors"
	"github.com/labstack/echo/v5"
)

type GuardWithSecret struct {
	secretKey string
}

func New(secretKey string) GuardWithSecret {
	return GuardWithSecret{
		secretKey: secretKey,
	}
}

func (m *GuardWithSecret) GuardWithSecretMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := c.Request()
		key := req.Header.Get("x-api-key")
		if key != m.secretKey {
			return customerrors.NewAppHTTPError(
				http.StatusForbidden,
				"Secret key not provided",
				errors.New("secret key not provided"),
			)
		}
		return next(c)
	}
}
