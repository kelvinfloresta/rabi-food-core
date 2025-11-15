package errs

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrProductNotFound = newErr("PRODUCT_NOT_FOUND", http.StatusNotFound)
)

func ProductNotFound(productIDs ...string) *AppError {
	if len(productIDs) == 0 {
		return ErrProductNotFound
	}

	e := *ErrProductNotFound
	e.Code = fmt.Sprintf("%s__%s", ErrProductNotFound.Code, strings.Join(productIDs, "_"))

	// To allow errors.Is checks
	e.err = ErrProductNotFound

	return &e
}
