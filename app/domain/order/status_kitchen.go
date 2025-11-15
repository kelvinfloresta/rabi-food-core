package order

import "slices"

// FulfillmentStatus tracks the preparation lifecycle inside the fulfillment.
type FulfillmentStatus string

const (
	FulfillmentPending      FulfillmentStatus = "pending"       // order created, not accepted yet
	FulfillmentConfirmed    FulfillmentStatus = "confirmed"     // accepted by the fulfillment
	FulfillmentInProduction FulfillmentStatus = "in_production" // being prepared
	FulfillmentReady        FulfillmentStatus = "ready"         // finished; ready for pickup/dispatch
)

var fulfillmentAllowed = map[FulfillmentStatus][]FulfillmentStatus{
	FulfillmentPending:      {FulfillmentConfirmed},
	FulfillmentConfirmed:    {FulfillmentInProduction},
	FulfillmentInProduction: {FulfillmentReady},
}

func (k FulfillmentStatus) CanTransition(to FulfillmentStatus) bool {
	return slices.Contains(fulfillmentAllowed[k], to)
}
