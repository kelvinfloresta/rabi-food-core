package errs

import (
	"net/http"
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
)
