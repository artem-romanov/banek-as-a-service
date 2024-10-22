package customvalidator

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	customerrors "baneks.com/internal/custom_errors"
	"github.com/go-playground/validator"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

var CustomValidator *validator.Validate = InitializeValidator()

func InitializeValidator() *validator.Validate {
	v := validator.New()
	// Allow validator read data from json annotation field.
	// This allows us to not send struct fields to the user.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(
			fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// initilize custom validation functions...
	inbetweenYearsValidator, err := CreateValidateMemeYear()
	if err != nil {
		log.Fatalf("can't create inbetweenYearsValidator: {%s}", err.Error())
	}
	v.RegisterValidation("is-correct-meme-year", inbetweenYearsValidator)

	return v
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

type validationFunc func(i interface{}) error

func ValidateRequest(checker validationFunc, data interface{}) *customerrors.AppHttpError {
	if checker == nil {
		return &customerrors.AppHttpError{
			Code:     http.StatusInternalServerError,
			Message:  "Can't validate your request",
			Internal: fmt.Errorf("checker in ValidateRequest is nil, did you forget to provide it?"),
		}
	}

	err := checker(data)
	if err == nil {
		return nil
	}

	if verr, ok := err.(validator.ValidationErrors); ok {
		return customerrors.NewAppHTTPError(http.StatusBadRequest, GetFancyErrors(&verr), fmt.Errorf("Request validation error: %s", err))
	}

	return customerrors.NewAppHTTPError(http.StatusBadRequest, err.Error(), err)
}
