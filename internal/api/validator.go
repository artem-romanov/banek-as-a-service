package server

import (
	customvalidator "baneks.com/internal/utils/validator"
	"github.com/go-playground/validator"
)

type AppValidator struct {
	validator *validator.Validate
}

func (cv *AppValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func CreateValidator() *AppValidator {
	v := customvalidator.CustomValidator
	return &AppValidator{
		validator: v,
	}
}
