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

// Error wraps an error with an HTTP status, application code, and message
type Error struct {
	HTTPCode int    `json:"-"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Err      error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s (code: %d) -> %v", e.Message, e.Code, e.Err)
	}
	return fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
}

// NewError creates a new application error
func NewError(httpCode, appCode int, message string, err error) *Error {
	return &Error{
		HTTPCode: httpCode,
		Code:     appCode,
		Message:  message,
		Err:      err,
	}
}

// -------------------------
// Error Handlers
// -------------------------

// BadRequest handles 400 - Bad Request consistently.
func BadRequest(ctx *gin.Context, err error) {
	// Case 1: Empty body
	if errors.Is(err, io.EOF) {
		ctx.Error(NewError(
			http.StatusBadRequest,
			constants.ErrCodeInvalidRequest,
			"request body is required but was empty",
			err,
		))
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

		ctx.Error(NewError(
			http.StatusBadRequest,
			constants.ErrCodeInvalidRequest,
			formatValidationErrors(errorsMap),
			err,
		))
		return
	}

	// Case 3: Other JSON/binding errors
	ctx.Error(NewError(
		http.StatusBadRequest,
		constants.ErrCodeInvalidRequest,
		"invalid request body",
		err,
	))
}

// Helper to convert validation errors map to string
func formatValidationErrors(errorsMap map[string]string) string {
	parts := make([]string, 0, len(errorsMap))
	for field, msg := range errorsMap {
		parts = append(parts, fmt.Sprintf("%s: %s", field, msg))
	}
	return strings.Join(parts, "; ")
}

// InternalServerError handles 500 - Internal Server Error
func InternalServerError(ctx *gin.Context, err error) {
	ctx.Error(NewError(
		http.StatusInternalServerError,
		constants.ErrCodeInternalServer,
		"internal server error",
		err,
	))
}

// NotFound handles 404 - Not Found
func NotFound(ctx *gin.Context, msg string) {
	ctx.Error(NewError(
		http.StatusNotFound,
		constants.ErrCodeUserNotFound,
		msg,
		nil,
	))
}

// PostgresError maps pg errors into AppError with proper codes
func PostgresError(ctx *gin.Context, err error) {
	// Handle sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Error(NewError(
			http.StatusNotFound,
			constants.ErrCodeUserNotFound,
			"resource not found",
			err,
		))
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			ctx.Error(NewError(http.StatusConflict, constants.ErrCodeUserAlreadyExists, "resource already exists", err))
		case "23503": // foreign_key_violation
			ctx.Error(NewError(http.StatusBadRequest, constants.ErrCodeInvalidRequest, "invalid reference, foreign key constraint failed", err))
		case "23502": // not_null_violation
			ctx.Error(NewError(http.StatusBadRequest, constants.ErrCodeInvalidRequest, "required field missing", err))
		case "23514": // check_violation
			ctx.Error(NewError(http.StatusBadRequest, constants.ErrCodeInvalidRequest, "check constraint failed", err))
		default:
			ctx.Error(NewError(http.StatusInternalServerError, constants.ErrCodeInternalServer, "database error", err))
		}
		return
	}

	// Fallback for all other errors
	InternalServerError(ctx, err)
}
