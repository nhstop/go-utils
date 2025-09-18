package apperr

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/busnosh/go-utils/pkg/constants"
	"github.com/jackc/pgx/v5/pgconn"
)

// PostgresError maps PostgreSQL errors into CodedError
func PostgresError(err error) *CodedError {
	// Handle sql.ErrNoRows
	if errors.Is(err, sql.ErrNoRows) {
		return NewError(ErrorParams{
			HTTPCode: http.StatusNotFound,
			Code:     constants.DBNotFound,
			Message:  "resource not found",
			Err:      err,
		})
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return NewError(ErrorParams{
				HTTPCode: http.StatusConflict,
				Code:     constants.DBUniqueViolation,
				Message:  "resource already exists",
				Err:      err,
			})
		case "23503": // foreign_key_violation
			return NewError(ErrorParams{
				HTTPCode: http.StatusBadRequest,
				Code:     constants.DBForeignKeyViolation,
				Message:  "invalid reference, foreign key constraint failed",
				Err:      err,
			})
		case "23502": // not_null_violation
			return NewError(ErrorParams{
				HTTPCode: http.StatusBadRequest,
				Code:     constants.InvalidRequest,
				Message:  "required field missing",
				Err:      err,
			})
		case "23514": // check_violation
			return NewError(ErrorParams{
				HTTPCode: http.StatusBadRequest,
				Code:     constants.InvalidRequest,
				Message:  "check constraint failed",
				Err:      err,
			})
		default: // all other Postgres errors
			return NewError(ErrorParams{
				HTTPCode: http.StatusInternalServerError,
				Code:     constants.DBError,
				Message:  "Application error",
				Err:      err,
			})
		}
	}

	// Fallback for all other errors
	return NewError(ErrorParams{
		HTTPCode: http.StatusInternalServerError,
		Code:     constants.InternalServer,
		Message:  "internal server error",
		Err:      err,
	})
}
