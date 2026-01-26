package server

import (
	"errors"
	"net/http"

	customerrors "baneks.com/internal/custom_errors"
	"github.com/labstack/echo/v5"
)

// CreateErrorHandler currently disabled
// TODO: think about removing it in future
func CreateErrorHandler() func(*echo.Context, error) {
	return func(c *echo.Context, err error) {
		if resp, err := echo.UnwrapResponse(c.Response()); err != nil {
			// response already been sent to the client
			if resp.Committed {
				return
			}
		}

		code := http.StatusInternalServerError
		var sc echo.HTTPStatusCoder
		if errors.As(err, &sc) {
			if tmp := sc.StatusCode(); tmp != 0 {
				code = tmp
			}
		}

		errorMap := make(map[string]interface{})

		var httpError *echo.HTTPError
		var appHttpError *customerrors.AppHttpError
		var errorText interface{} = "Internal server error"
		switch {
		case errors.As(err, &httpError):
			code = httpError.Code
			errorMap["message"] = httpError.Message
			errorText = httpError.Message
		case errors.As(err, &appHttpError):
			code = appHttpError.Code
			errorMap["message"] = appHttpError.Message
			if appHttpError.Internal != nil {
				errorText = appHttpError.Internal.Error()
			} else {
				errorText = appHttpError.Message
			}
		default:
			errorMap["message"] = http.StatusText(code)
			errorText = err.Error()
		}

		if code >= 500 {
			c.Logger().Error("internal error", "err", errorText)
		}
		c.JSON(code, errorMap)
	}
}
