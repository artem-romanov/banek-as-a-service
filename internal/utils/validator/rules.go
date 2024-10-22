package customvalidator

import (
	"errors"
	"time"

	"github.com/go-playground/validator"
)

// / CUSTOM VALIDATION RULES

func CreateValidateYearInbetween(yearMin int, yearMax int) (func(fl validator.FieldLevel) bool, error) {
	if yearMin > yearMax {
		return nil, errors.New("minYear can't be higher that yearMax")
	}

	return func(fl validator.FieldLevel) bool {
		can := fl.Field().CanInt()
		if !can {
			return false
		}
		field := int(fl.Field().Int())
		if field >= yearMin && field <= yearMax {
			return true
		}
		return false
	}, nil
}

func CreateValidateYearInbetweenNow(yearMin int) (func(fl validator.FieldLevel) bool, error) {
	now := time.Now()
	currentYear := now.Year()
	validator, error := CreateValidateYearInbetween(yearMin, currentYear)
	if error != nil {
		return nil, error
	}
	return validator, nil
}

func CreateValidateMemeYear() (func(fl validator.FieldLevel) bool, error) {
	// 2015 is the last available year available on the website
	// https://idiod.qabyldau.com/random/
	return CreateValidateYearInbetweenNow(2015)
}
