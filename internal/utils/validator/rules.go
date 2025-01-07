package customvalidator

import (
	"errors"
	"time"

	"github.com/go-playground/validator"
)

// / CUSTOM VALIDATION RULES

func CreateValidateYearInbetween(yearMin int, yearMax func() int) (func(fl validator.FieldLevel) bool, error) {
	if yearMin > yearMax() {
		return nil, errors.New("minYear can't be higher that yearMax")
	}

	return func(fl validator.FieldLevel) bool {
		can := fl.Field().CanInt()
		if !can {
			return false
		}
		field := int(fl.Field().Int())
		if field >= yearMin && field <= yearMax() {
			return true
		}
		return false
	}, nil
}

func CreateValidateYearInbetweenNow(yearMin int) (func(fl validator.FieldLevel) bool, error) {
	now := time.Now()
	validator, error := CreateValidateYearInbetween(yearMin, now.Year)
	if error != nil {
		return nil, error
	}
	return validator, nil
}

func CreateValidateMemeYear() (func(fl validator.FieldLevel) bool, error) {
	// 2015 is the last available year on the website
	// https://idiod.qabyldau.com/random/
	return CreateValidateYearInbetweenNow(2015)
}
