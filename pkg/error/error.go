package apperr

import (
	"fmt"
	"net/http"

	"github.com/busnosh/go-utils/pkg/constants"
)

// -------------------------
// CodedError & Constructor
// -------------------------

type CodedError struct {
	HTTPCode int    // HTTP status code
	Code     int    // Application error code
	Message  string // Human-readable message
	Err      error  // Underlying error
}

func (c *CodedError) Error() string {
	if c.Err != nil {
		return fmt.Sprintf("%s (code: %d, http: %d) -> %v", c.Message, c.Code, c.HTTPCode, c.Err)
	}
	return fmt.Sprintf("%s (code: %d, http: %d)", c.Message, c.Code, c.HTTPCode)
}

// Optional parameters struct for NewError
type ErrorParams struct {
	HTTPCode int
	Code     int
	Message  string
	Err      error
}

// NewError creates a new CodedError with defaults and optional overrides
func NewError(params ErrorParams) *CodedError {
	e := &CodedError{
		HTTPCode: http.StatusInternalServerError,
		Code:     constants.InternalServer,
		Message:  "internal server error",
		Err:      nil,
	}

	if params.HTTPCode != 0 {
		e.HTTPCode = params.HTTPCode
	}
	if params.Code != 0 {
		e.Code = params.Code
	}
	if params.Message != "" {
		e.Message = params.Message
	}
	if params.Err != nil {
		e.Err = params.Err
	}

	return e
}

// -------------------------
// Error Handlers
// -------------------------

// Helper to format validation errors map

// InternalServerError handles 500 errors
func InternalServerError(err error) *CodedError {
	msg := "internal server error"
	if err != nil && err.Error() != "" {
		msg = err.Error()
	}

	return NewError(ErrorParams{
		HTTPCode: http.StatusInternalServerError,
		Message:  msg,
		Err:      err,
	})
}

// NotFound handles 404 errors
func NotFound(msg string) *CodedError {
	return NewError(ErrorParams{
		HTTPCode: http.StatusNotFound,
		Message:  msg,
	})
}
