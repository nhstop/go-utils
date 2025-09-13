package err

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Error wraps an error with code & message
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

func NewError(code int, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// ----------------- Helpers -----------------

func BadRequest(ctx *gin.Context, err error) {
	if errors.Is(err, io.EOF) {
		ctx.Error(NewError(http.StatusBadRequest, "request body is required", err))
		return
	}

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrs {
			msg := fieldErr.Field() + " failed on '" + fieldErr.Tag() + "'"
			if fieldErr.Param() != "" {
				msg += " (param: " + fieldErr.Param() + ")"
			}
			errorsMap[fieldErr.Field()] = msg
		}
		ctx.Error(NewError(http.StatusBadRequest, formatValidationErrors(errorsMap), err))
		return
	}

	ctx.Error(NewError(http.StatusBadRequest, "invalid request body", err))
}

func InternalServerError(ctx *gin.Context, err error) {
	ctx.Error(NewError(http.StatusInternalServerError, err.Error(), err))
}

func NotFound(ctx *gin.Context, msg string) {
	ctx.Error(NewError(http.StatusNotFound, msg, nil))
}

// func PostgresError(ctx *gin.Context, err error) {
// 	if errors.Is(err, sql.ErrNoRows) {
// 		ctx.Error(NewError(http.StatusNotFound, "Resource not found", err))
// 		return
// 	}

// 	var pgErr *pgconn.PgError
// 	if errors.As(err, &pgErr) {
// 		switch pgErr.Code {
// 		case "23505":
// 			ctx.Error(NewError(http.StatusConflict, "Resource already exists", err))
// 		case "23503", "23502", "23514":
// 			ctx.Error(NewError(http.StatusBadRequest, "Database constraint failed", err))
// 		default:
// 			ctx.Error(NewError(http.StatusInternalServerError, "Database error", err))
// 		}
// 		return
// 	}

// 	InternalServerError(ctx, err)
// }

// format validation map to string
func formatValidationErrors(errorsMap map[string]string) string {
	result := ""
	for field, msg := range errorsMap {
		result += field + ": " + msg + "; "
	}
	return result
}
