package util

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func GetFancyErrors(validationErrors *validator.ValidationErrors) []ValidationError {
	errors := []ValidationError{}
	for _, validationError := range *validationErrors {
		fieldName := validationError.Field()
		reason := validationError.ActualTag()
		errors = append(errors, ValidationError{
			Field:  fieldName,
			Reason: "failed on " + reason + " validation rule",
		})
	}
	return errors
}

func Validate(c echo.Context, data interface{}) *echo.HTTPError {
	err := c.Validate(data)
	if err == nil {
		return nil
	}

	if verr, ok := err.(validator.ValidationErrors); ok {
		return echo.NewHTTPError(http.StatusBadRequest, GetFancyErrors(&verr))
	}

	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}
