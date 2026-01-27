package middlewares

import (
	"errors"
	"net/http"

	customerrors "baneks.com/internal/custom_errors"
	"github.com/labstack/echo/v5"
)

type guardWithSecret struct {
	secretKey string
}

func NewGuardWithSecret(secretKey string) guardWithSecret {
	return guardWithSecret{
		secretKey: secretKey,
	}
}

func (m *guardWithSecret) GuardWithSecretMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
