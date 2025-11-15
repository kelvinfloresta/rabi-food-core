package order

import "slices"

// DeliveryStatus tracks the last-mile logistics lifecycle.
// If the order is pickup-at-store only, it can remain "pending" and jump to "delivered".
type DeliveryStatus string

const (
	DeliveryPending   DeliveryStatus = "pending"          // awaiting fulfillment=ready
	DeliveryAssigned  DeliveryStatus = "assigned"         // driver assigned / courier matched
	DeliveryPickedUp  DeliveryStatus = "picked_up"        // courier picked up at store
	DeliveryOutFor    DeliveryStatus = "out_for_delivery" // on the way
	DeliveryDelivered DeliveryStatus = "delivered"        // customer received
	DeliveryCancelled DeliveryStatus = "cancelled"        // cancelled before dispatch
	DeliveryRefunded  DeliveryStatus = "refunded"         // refunded after delivery/charge
)

var deliveryAllowed = map[DeliveryStatus][]DeliveryStatus{
	DeliveryPending:   {DeliveryAssigned, DeliveryCancelled},
	DeliveryAssigned:  {DeliveryPickedUp, DeliveryCancelled},
	DeliveryPickedUp:  {DeliveryOutFor, DeliveryCancelled},
	DeliveryOutFor:    {DeliveryDelivered, DeliveryRefunded},
	DeliveryDelivered: {DeliveryRefunded},
}

func (from DeliveryStatus) CanTransition(to DeliveryStatus) bool {
	return slices.Contains(deliveryAllowed[from], to)
}
