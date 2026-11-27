package errs

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNoValuesToUpdate = newErr("NO_VALUES_TO_UPDATE", http.StatusBadRequest)
	ErrConflict         = newErr("CONFLICT", http.StatusConflict)
)

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == "23505"
	}

	return false
}

// DatabaseErrorInfo reports whether the database caused err and, when the database answered, its SQLSTATE. It
// never returns the driver's detail message, which may hold the values of the failed row.
func DatabaseErrorInfo(err error) (string, bool) {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code, true
	}

	if _, ok := errors.AsType[*pgconn.ConnectError](err); ok {
		return "", true
	}

	return "", false
}
