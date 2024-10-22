package customerrors

import (
	"fmt"
	"net/http"
)

// Yeah, it looks the same as echo.HTTPError, but hear me out.
// *Here* "Internal" is an error which actually happend in the app.
// Idea is to return AppHttpError instead of HTTPError from controllers.
// Send Message to user *and* log Internal error with all nifty things such as request uri.
type AppHttpError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"`
}

func (e *AppHttpError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("code=%d, message=%v, internal=%v", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("code=%d, message=%v", e.Code, e.Message)
}

func NewAppHTTPError(code int, message interface{}, err error) *AppHttpError {
	return &AppHttpError{
		Code:     code,
		Message:  message,
		Internal: err,
	}
}

func NewAppBindError(err error) *AppHttpError {
	return &AppHttpError{
		Code:     http.StatusInternalServerError,
		Message:  "Unable to parse data",
		Internal: err,
	}
}
