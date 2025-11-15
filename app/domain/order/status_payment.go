package order

import "slices"

// PaymentStatus tracks the financial lifecycle, independent from production/delivery.
type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "pending"    // not initiated or waiting user action
	PaymentAuthorized PaymentStatus = "authorized" // auth ok, yet to capture
	PaymentPaid       PaymentStatus = "paid"       // captured/settled
	PaymentFailed     PaymentStatus = "failed"     // attempt failed
	PaymentCancelled  PaymentStatus = "cancelled"  // voided before capture
	PaymentRefunded   PaymentStatus = "refunded"   // refunded after capture
)

var paymentAllowed = map[PaymentStatus][]PaymentStatus{
	PaymentPending:    {PaymentAuthorized, PaymentPaid, PaymentFailed, PaymentCancelled},
	PaymentAuthorized: {PaymentPaid, PaymentCancelled, PaymentFailed},
	PaymentPaid:       {PaymentRefunded},
	PaymentFailed:     {}, // terminal; new attempt = new intent
	PaymentCancelled:  {}, // terminal for this intent
	PaymentRefunded:   {}, // terminal
}

func (from PaymentStatus) CanTransition(to PaymentStatus) bool {
	return slices.Contains(paymentAllowed[from], to)
}
