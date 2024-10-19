package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateErrorHandler() func(error, echo.Context) {
	return func(err error, c echo.Context) {
		c.Logger().Debug(c.Request().URL)
		code := http.StatusInternalServerError
		he, ok := err.(*echo.HTTPError)
		if ok {
			code = he.Code
		}

		c.Logger().Error(err)
		errorMap := make(map[string]interface{})
		errorMap["message"] = he.Message

		c.JSON(code, errorMap)
	}
}
