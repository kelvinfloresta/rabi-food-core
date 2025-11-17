package order_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/domain/order"
	g "rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/errs"
)

type PatchFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

func (c *OrderCase) Patch(
	ctx context.Context,
	filter g.PatchFilter,
	values g.PatchValues,
) (bool, error) {
	session := app_context.GetSession(ctx)
	if !session.Role.IsBackoffice() {
		filter.TenantID = session.TenantID
	}

	if session.Role.IsUser() {
		return false, errs.ErrForbidden
	}

	deliveryStatusIn := []order.DeliveryStatus(nil)
	if values.DeliveryStatus != "" {
		deliveryStatusIn = values.DeliveryStatus.GetPrerequisites()
		if len(deliveryStatusIn) == 0 {
			return false, errs.StatusNotModifiable(string(values.DeliveryStatus))
		}
	}

	fulfillmentStatusIn := []order.FulfillmentStatus(nil)
	if values.FulfillmentStatus != "" {
		fulfillmentStatusIn = values.FulfillmentStatus.GetPrerequisites()
		if len(fulfillmentStatusIn) == 0 {
			return false, errs.StatusNotModifiable(string(values.FulfillmentStatus))
		}
	}

	paymentStatusIn := []order.PaymentStatus(nil)
	if values.PaymentStatus != "" {
		paymentStatusIn = values.PaymentStatus.GetPrerequisites()
		if len(paymentStatusIn) == 0 {
			return false, errs.StatusNotModifiable(string(values.PaymentStatus))
		}
	}

	return c.gateway.Patch(g.PatchFilter{
		ID:                  filter.ID,
		TenantID:            filter.TenantID,
		DeliveryStatusIn:    deliveryStatusIn,
		FulfillmentStatusIn: fulfillmentStatusIn,
		PaymentStatusIn:     paymentStatusIn,
	}, values)
}
