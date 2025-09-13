package apperr

import (
	"database/sql"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

// Error wraps an error with an HTTP status code and message
type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewError creates a new Error
func NewError(code int, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// -------------------------
// Error Handlers
// -------------------------

// BadRequest handles 400 - Bad Request
func BadRequest(ctx *gin.Context, err error) {
	// Case 1: Empty body
	if errors.Is(err, io.EOF) {
		ctx.Error(NewError(http.StatusBadRequest, "request body is required but was empty", err))
		return
	}

	// Case 2: Validation errors
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrs {
			msg := fieldErr.Field() + " failed on '" + fieldErr.Tag() + "'"
			if fieldErr.Param() != "" {
				msg += " (param: " + fieldErr.Param() + ")"
			}
			errorsMap[fieldErr.Field()] = msg
		}

		// Attach Error to context, store validation map in Err
		ctx.Error(NewError(http.StatusBadRequest, formatValidationErrors(errorsMap), err))
		return
	}

	// Case 3: Other JSON/binding errors
	ctx.Error(NewError(http.StatusBadRequest, "invalid request body", err))
}

// Helper to convert validation errors map to string
func formatValidationErrors(errorsMap map[string]string) string {
	result := ""
	for field, msg := range errorsMap {
		result += field + ": " + msg + "; "
	}
	return result
}

// InternalServerError handles 500 - Internal Server Error
func InternalServerError(ctx *gin.Context, err error) {
	ctx.Error(NewError(http.StatusInternalServerError, err.Error(), err))
}

// NotFound handles 404 - Not Found
func NotFound(ctx *gin.Context, msg string) {
	ctx.Error(NewError(http.StatusNotFound, msg, nil))
}
func PostgresError(ctx *gin.Context, err error) {
	// Handle sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Error(NewError(http.StatusNotFound, "Resource not found", err))
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {

		switch pgErr.Code {
		case "23505": // unique_violation
			ctx.Error(NewError(http.StatusConflict, "Resource already exists", err))
		case "23503": // foreign_key_violation
			ctx.Error(NewError(http.StatusBadRequest, "Invalid reference, foreign key constraint failed", err))
		case "23502": // not_null_violation
			ctx.Error(NewError(http.StatusBadRequest, "Required field missing", err))
		case "23514": // check_violation
			ctx.Error(NewError(http.StatusBadRequest, "Check constraint failed", err))
		default:
			ctx.Error(NewError(http.StatusInternalServerError, "Database error", err))
		}
		return
	}

	// Fallback for all other errors
	InternalServerError(ctx, err)
}
