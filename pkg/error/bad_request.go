package apperr

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/busnosh/go-utils/pkg/constants"
	"github.com/go-playground/validator/v10"
)

// BadRequest handles 400 - Bad Request consistently
func BadRequest(err error) *CodedError {
	// Case 1: Empty body
	if errors.Is(err, io.EOF) {
		return NewError(ErrorParams{
			HTTPCode: http.StatusBadRequest,
			Code:     constants.InvalidRequest,
			Message:  "request body is required but was empty",
			Err:      err,
		})

	}

	// Case 2: Validation errors
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrs {
			var msg string
			switch fieldErr.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", strings.ToLower(fieldErr.Field()))
			case "email":
				msg = fmt.Sprintf("%s must be a valid email", strings.ToLower(fieldErr.Field()))
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", strings.ToLower(fieldErr.Field()), fieldErr.Param())
			case "max":
				msg = fmt.Sprintf("%s cannot be longer than %s characters", strings.ToLower(fieldErr.Field()), fieldErr.Param())
			default:
				msg = fmt.Sprintf("%s failed on '%s'", strings.ToLower(fieldErr.Field()), fieldErr.Tag())
				if fieldErr.Param() != "" {
					msg += fmt.Sprintf(" (param: %s)", fieldErr.Param())
				}
			}
			errorsMap[strings.ToLower(fieldErr.Field())] = msg
		}

		return NewError(ErrorParams{
			HTTPCode: http.StatusBadRequest,
			Code:     constants.InvalidRequest,
			Message:  formatValidationErrors(errorsMap),
			Err:      err,
		})

	}

	// Case 3: Other JSON/binding errors
	return NewError(ErrorParams{
		HTTPCode: http.StatusBadRequest,
		Code:     constants.InvalidRequest,
		Message:  "invalid request body",
		Err:      err,
	})
}

func formatValidationErrors(errorsMap map[string]string) string {
	parts := make([]string, 0, len(errorsMap))
	for field, msg := range errorsMap {
		parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
	}
	return strings.Join(parts, "; ")
}
