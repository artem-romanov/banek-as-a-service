package customerrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AppHttpError is re-implementation of original echo.HTTPError
//
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

func (e *AppHttpError) StatusCode() int {
	return e.Code
}

func (e *AppHttpError) MessageString() string {
	switch v := e.Message.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return ""
	}
}

// MarshalJSON reimplements original marshal to support
// inside echo DefaultHttpHandler.
//
// Because to print message we either need to use original HttpError (which doesn't have internal error),
// or implement json.Marshaler
func (e *AppHttpError) MarshalJSON() ([]byte, error) {
	type Alias AppHttpError
	return json.Marshal(&struct {
		*Alias
		Message string `json:"message"`
	}{
		Alias:   (*Alias)(e),
		Message: e.MessageString(),
	})
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
