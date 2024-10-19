package server

import (
	"reflect"
	"strings"

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
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(
			fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &AppValidator{
		validator: v,
	}
}
