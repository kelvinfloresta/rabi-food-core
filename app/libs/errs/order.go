package errs

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrInvalidFulfillmentTransition = newErr("ORDER_KITCHEN_INVALID_TRANSITION", http.StatusConflict)
	ErrInvalidDeliveryTransition    = newErr("ORDER_DELIVERY_INVALID_TRANSITION", http.StatusConflict)
	ErrInvalidPaymentTransition     = newErr("ORDER_PAYMENT_INVALID_TRANSITION", http.StatusConflict)
	ErrPaymentNotCleared            = newErr("ORDER_PAYMENT_NOT_CLEARED", http.StatusPaymentRequired)
	ErrFulfillmentNotReady          = newErr("ORDER_KITCHEN_NOT_READY", http.StatusPreconditionFailed)
	ErrCannotDeliver                = newErr("ORDER_CANNOT_DELIVER", http.StatusPreconditionFailed)

	ErrProductNotFound = newErr("PRODUCT_NOT_FOUND", http.StatusNotFound)
	ErrOrderInvalid    = newErr("ORDER_INVALID", http.StatusBadRequest)

	ErrMissingParameter = newErr("MISSING_PARAMETER", http.StatusBadRequest)
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
