package errs

import (
	"fmt"
	"net/http"
)

// ErrInternal is the code sent for any error that is not an ApiError.
var ErrInternal = newErr("INTERNAL_ERROR", http.StatusInternalServerError)

// ApiError is sent to the client as {"code", "fullCode", "status"}. Code is a fixed value per error, for
// clients to branch on and for the error_code log field. FullCode adds the variable detail (IDs, statuses)
// for developers to read, and is omitted when the error has none.
type ApiError struct {
	Code     string `json:"code"`
	FullCode string `json:"fullCode,omitempty"`
	Status   int    `json:"status"`
	err      error  `json:"-"`
}

func newErr(code string, status int) *ApiError {
	return &ApiError{Code: code, Status: status}
}

func (e *ApiError) Error() string {
	code := e.fullCode()

	if e.err != nil {
		return fmt.Sprintf("%s: %v", code, e.err)
	}

	return code
}

func (e *ApiError) Unwrap() error {
	return e.err
}

func (e *ApiError) fullCode() string {
	if e.FullCode != "" {
		return e.FullCode
	}

	return e.Code
}
