package errs

import (
	"fmt"
)

type AppError struct {
	Code   string `json:"code"`
	Status int    `json:"status"`
	Err    error  `json:"-"`
}

func newErr(code string, status int) *AppError {
	return &AppError{Code: code, Status: status}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}

	return e.Code
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Wrap(cause error) *AppError {
	if e == nil {
		return nil
	}

	return &AppError{
		Code:   e.Code,
		Status: e.Status,
		Err:    cause,
	}
}

func (e *AppError) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `{"code":"%s","status":%d}`, e.Code, e.Status), nil
}
