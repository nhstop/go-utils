package apperr

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/busnosh/go-utils/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
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
		Code:     constants.ErrCodeInternalServer,
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

// BadRequest handles 400 - Bad Request consistently
func BadRequest(ctx *gin.Context, err error) {
	// Case 1: Empty body
	if errors.Is(err, io.EOF) {
		ctx.Error(NewError(ErrorParams{
			HTTPCode: http.StatusBadRequest,
			Code:     constants.ErrCodeInvalidRequest,
			Message:  "request body is required but was empty",
			Err:      err,
		}))
		return
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

		ctx.Error(NewError(ErrorParams{
			HTTPCode: http.StatusBadRequest,
			Code:     constants.ErrCodeInvalidRequest,
			Message:  formatValidationErrors(errorsMap),
			Err:      err,
		}))
		return
	}

	// Case 3: Other JSON/binding errors
	ctx.Error(NewError(ErrorParams{
		HTTPCode: http.StatusBadRequest,
		Code:     constants.ErrCodeInvalidRequest,
		Message:  "invalid request body",
		Err:      err,
	}))
}

// Helper to format validation errors map
func formatValidationErrors(errorsMap map[string]string) string {
	parts := make([]string, 0, len(errorsMap))
	for field, msg := range errorsMap {
		parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
	}
	return strings.Join(parts, "; ")
}

// InternalServerError handles 500 errors
func InternalServerError(ctx *gin.Context, err error) {
	msg := "internal server error"
	if err != nil && err.Error() != "" {
		msg = err.Error()
	}

	ctx.Error(NewError(ErrorParams{
		HTTPCode: http.StatusInternalServerError,
		Code:     constants.ErrCodeInternalServer,
		Message:  msg,
		Err:      err,
	}))
}

// NotFound handles 404 errors
func NotFound(ctx *gin.Context, msg string) {
	ctx.Error(NewError(ErrorParams{
		HTTPCode: http.StatusNotFound,
		Code:     constants.ErrCodeUserNotFound,
		Message:  msg,
	}))
}

// PostgresError maps PostgreSQL errors into CodedError
func PostgresError(ctx *gin.Context, err error) {
	// Handle sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Error(NewError(ErrorParams{
			HTTPCode: http.StatusNotFound,
			Code:     constants.ErrCodeUserNotFound,
			Message:  "resource not found",
			Err:      err,
		}))
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			ctx.Error(NewError(ErrorParams{
				HTTPCode: http.StatusConflict,
				Code:     constants.ErrCodeUserAlreadyExists,
				Message:  "resource already exists",
				Err:      err,
			}))
		case "23503": // foreign_key_violation
			ctx.Error(NewError(ErrorParams{
				HTTPCode: http.StatusBadRequest,
				Code:     constants.ErrCodeInvalidRequest,
				Message:  "invalid reference, foreign key constraint failed",
				Err:      err,
			}))
		case "23502": // not_null_violation
			ctx.Error(NewError(ErrorParams{
				HTTPCode: http.StatusBadRequest,
				Code:     constants.ErrCodeInvalidRequest,
				Message:  "required field missing",
				Err:      err,
			}))
		case "23514": // check_violation
			ctx.Error(NewError(ErrorParams{
				HTTPCode: http.StatusBadRequest,
				Code:     constants.ErrCodeInvalidRequest,
				Message:  "check constraint failed",
				Err:      err,
			}))
		default:
			ctx.Error(NewError(ErrorParams{
				HTTPCode: http.StatusInternalServerError,
				Code:     constants.ErrCodeInternalServer,
				Message:  "database error",
				Err:      err,
			}))
		}
		return
	}

	// Fallback for all other errors
	InternalServerError(ctx, err)
}
