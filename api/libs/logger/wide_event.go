package logger

import (
	"context"
	"errors"
	"sync"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/rs/zerolog"
)

const wideEventKey = "wideEventKey"

// WideEvent is a struct that captures all relevant information about an HTTP request and its context for logging purposes.
// It is designed to be used in the Logging middleware to emit structured logs with rich contextual data.
type WideEvent struct {
	mu sync.Mutex

	// Event is a custom field to describe the specific event or action being logged
	Event string

	Method     string
	Path       string
	Query      string
	StatusCode int
	LatencyMs  int64
	RequestID  string

	orderID        string
	categoryID     string
	productID      string
	subscriptionID string
	// userID is who made the request (the session user), set by the session middleware.
	userID string
	// selectedTenantID is the tenant a customer selected for the request. The tenant a member acts for is not
	// logged: it follows from userID.
	selectedTenantID string

	// Payment
	externalPaymentID         string
	originalExternalPaymentID string
	paymentIdempotencyHit     *bool
	paymentProvider           string
	paymentAmountCents        *int64
	paymentChargeType         string
	paymentStatus             string
	paymentRequiresAction     *bool

	// Delivery Scheduling
	deliverySchedulingDate     string
	deliverySchedulingWeekday  string
	deliverySchedulingTimezone string
	deliverySchedulingCount    *int

	err       string
	errorCode string
}

// --- Usecase setters (thread-safe) ---

func (w *WideEvent) SetOrderID(id string) {
	w.mu.Lock()
	w.orderID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetCategoryID(id string) {
	w.mu.Lock()
	w.categoryID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetProductID(id string) {
	w.mu.Lock()
	w.productID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetSubscriptionID(id string) {
	w.mu.Lock()
	w.subscriptionID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetUserID(id string) {
	w.mu.Lock()
	w.userID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetSelectedTenantID(id string) {
	w.mu.Lock()
	w.selectedTenantID = id
	w.mu.Unlock()
}

func (w *WideEvent) TrackOrderPayment(orderID, externalPaymentID string) {
	w.mu.Lock()
	w.orderID = orderID
	w.externalPaymentID = externalPaymentID
	w.mu.Unlock()
}

func (w *WideEvent) TrackSubscriptionPayment(subscriptionID, externalPaymentID, provider string, amountCents int64) {
	w.mu.Lock()
	w.subscriptionID = subscriptionID
	w.externalPaymentID = externalPaymentID
	w.paymentProvider = provider
	w.paymentAmountCents = &amountCents
	w.mu.Unlock()
}

func (w *WideEvent) MarkPaymentIdempotent() {
	w.mu.Lock()
	v := true
	w.paymentIdempotencyHit = &v
	w.mu.Unlock()
}

func (w *WideEvent) TrackCharge(chargeType, status string, requiresAction bool) {
	w.mu.Lock()
	w.paymentChargeType = chargeType
	w.paymentStatus = status
	w.paymentRequiresAction = &requiresAction
	w.mu.Unlock()
}

func (w *WideEvent) TrackPaymentConflict(originalExternalPaymentID string) {
	w.mu.Lock()
	w.originalExternalPaymentID = originalExternalPaymentID
	w.mu.Unlock()
}

func (w *WideEvent) TrackDeliveryScheduling(date, weekday, timezone string, count int) {
	w.mu.Lock()
	w.deliverySchedulingDate = date
	w.deliverySchedulingWeekday = weekday
	w.deliverySchedulingTimezone = timezone
	w.deliverySchedulingCount = &count
	w.mu.Unlock()
}

func (w *WideEvent) SetErrorCode(code string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.errorCode = code
}

// FinishRequest records the request outcome. errorCode is the code sent to the client, logged even when
// err is nil (validation errors are not logged as errors, but their code is).
func (w *WideEvent) FinishRequest(statusCode int, errorCode string, err error, latencyMs int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.StatusCode = statusCode
	w.LatencyMs = latencyMs
	w.errorCode = errorCode
	if err != nil {
		w.err = err.Error()
	}
}

func (w *WideEvent) FinishCronJob(err error, latencyMs int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.LatencyMs = latencyMs
	if err != nil {
		w.err = err.Error()
		w.errorCode = errs.ErrInternal.Code
		if apiError, ok := errors.AsType[*errs.ApiError](err); ok {
			w.errorCode = apiError.Code
		}
	}
}

// reset clears all fields of the WideEvent, preparing it for reuse from the pool.
func (w *WideEvent) reset() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Event = ""
	w.Method = ""
	w.Path = ""
	w.Query = ""
	w.StatusCode = 0
	w.LatencyMs = 0
	w.RequestID = ""

	w.orderID = ""
	w.categoryID = ""
	w.productID = ""
	w.subscriptionID = ""
	w.userID = ""
	w.selectedTenantID = ""

	w.externalPaymentID = ""
	w.originalExternalPaymentID = ""
	w.paymentIdempotencyHit = nil
	w.paymentProvider = ""
	w.paymentAmountCents = nil
	w.paymentChargeType = ""
	w.paymentStatus = ""
	w.paymentRequiresAction = nil

	w.deliverySchedulingDate = ""
	w.deliverySchedulingWeekday = ""
	w.deliverySchedulingTimezone = ""
	w.deliverySchedulingCount = nil

	w.err = ""
	w.errorCode = ""
}

func (w *WideEvent) MarshalZerologObject(e *zerolog.Event) {
	w.mu.Lock()
	defer w.mu.Unlock()

	e.
		Str("method", w.Method).
		Str("path", w.Path).
		Str("query", w.Query).
		Int("status", w.StatusCode).
		Int64("latency_ms", w.LatencyMs)

	// request_id is also sent to the client (X-Request-ID header), so a reported failure can be found.
	if w.RequestID != "" {
		e.Str("request_id", w.RequestID)
	}

	if w.Event != "" {
		e.Str("event", w.Event)
	}

	if w.orderID != "" {
		e.Str("order_id", w.orderID)
	}
	if w.categoryID != "" {
		e.Str("category_id", w.categoryID)
	}
	if w.productID != "" {
		e.Str("product_id", w.productID)
	}
	if w.subscriptionID != "" {
		e.Str("subscription_id", w.subscriptionID)
	}
	if w.userID != "" {
		e.Str("user_id", w.userID)
	}
	if w.selectedTenantID != "" {
		e.Str("selected_tenant_id", w.selectedTenantID)
	}

	if w.externalPaymentID != "" {
		e.Str("payment.external_id", w.externalPaymentID)
	}
	if w.originalExternalPaymentID != "" {
		e.Str("payment.original_external_id", w.originalExternalPaymentID)
	}
	if w.paymentIdempotencyHit != nil {
		e.Bool("payment.idempotency_hit", *w.paymentIdempotencyHit)
	}
	if w.paymentProvider != "" {
		e.Str("payment.provider", w.paymentProvider)
	}
	if w.paymentAmountCents != nil {
		e.Int64("payment.amount_cents", *w.paymentAmountCents)
	}
	if w.paymentChargeType != "" {
		e.Str("payment.charge_type", w.paymentChargeType)
	}
	if w.paymentStatus != "" {
		e.Str("payment.status", w.paymentStatus)
	}
	if w.paymentRequiresAction != nil {
		e.Bool("payment.requires_action", *w.paymentRequiresAction)
	}

	if w.deliverySchedulingDate != "" {
		e.Str("delivery_scheduling.date", w.deliverySchedulingDate)
		e.Str("delivery_scheduling.weekday", w.deliverySchedulingWeekday)
		e.Str("delivery_scheduling.timezone", w.deliverySchedulingTimezone)
		if w.deliverySchedulingCount != nil {
			e.Int("delivery_scheduling.count", *w.deliverySchedulingCount)
		}
	}

	if w.err != "" {
		e.Str("error", w.err)
	}
	if w.errorCode != "" {
		e.Str("error_code", w.errorCode)
	}
}

var wideEventPool = sync.Pool{
	New: func() any {
		return &WideEvent{}
	},
}

// Acquire retrieves a WideEvent from the pool, resetting its fields before returning it for use.
func Acquire() *WideEvent {
	wideEvent := wideEventPool.Get().(*WideEvent)

	return wideEvent
}

// Release returns a WideEvent to the pool for reuse.
func Release(w *WideEvent) {
	w.reset()
	wideEventPool.Put(w)
}

func GetWideEvent(ctx context.Context) *WideEvent {
	if we, ok := ctx.Value(wideEventKey).(*WideEvent); ok && we != nil {
		return we
	}

	l := Get(ctx)
	l.Error().Msg("WideEvent not found in context")

	return &WideEvent{}
}

func WithWideEvent(ctx context.Context, we *WideEvent) context.Context {
	return context.WithValue(ctx, wideEventKey, we)
}
