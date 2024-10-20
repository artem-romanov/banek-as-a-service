package server

import (
	"errors"
	"net/http"

	customerrors "baneks.com/internal/custom_errors"
	"github.com/labstack/echo/v4"
)

func CreateErrorHandler() func(error, echo.Context) {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
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
			errorMap["message"] = err.Error()
			errorText = err.Error()
		}

		c.Logger().Error(errorText)
		c.JSON(code, errorMap)
	}
}
