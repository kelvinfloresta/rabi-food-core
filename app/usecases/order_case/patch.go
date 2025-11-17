package order_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/domain/order"
	g "rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/errs"
	"rabi-food-core/libs/logger"
)

type PatchFilter struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
}

func (c *OrderCase) Patch(
	ctx context.Context,
	filter PatchFilter,
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
			return false, errs.StatusNotModifiable(values.DeliveryStatus)
		}
	}

	fulfillmentStatusIn := []order.FulfillmentStatus(nil)
	if values.FulfillmentStatus != "" {
		fulfillmentStatusIn = values.FulfillmentStatus.GetPrerequisites()
		if len(fulfillmentStatusIn) == 0 {
			return false, errs.StatusNotModifiable(values.FulfillmentStatus)
		}
	}

	paymentStatusIn := []order.PaymentStatus(nil)
	if values.PaymentStatus != "" {
		paymentStatusIn = values.PaymentStatus.GetPrerequisites()
		if len(paymentStatusIn) == 0 {
			return false, errs.StatusNotModifiable(values.PaymentStatus)
		}
	}

	patched, err := c.gateway.Patch(g.PatchFilter{
		ID:                  filter.ID,
		TenantID:            filter.TenantID,
		DeliveryStatusIn:    deliveryStatusIn,
		FulfillmentStatusIn: fulfillmentStatusIn,
		PaymentStatusIn:     paymentStatusIn,
	}, values)

	if err != nil {
		return false, err
	}

	if patched {
		return true, nil
	}

	logger.Get(ctx).Warn().Str("orderID", filter.ID).Msg("order not patched")
	return c.handleNotPatched(ctx, filter, values)
}

func (c *OrderCase) handleNotPatched(ctx context.Context, filter PatchFilter, values g.PatchValues) (bool, error) {
	l := logger.Get(ctx).Warn().Str("orderID", filter.ID)

	l.Msg("checking if order exists for not patched case")
	orderFound, err := c.gateway.GetByID(g.GetByIDFilter{
		ID:       filter.ID,
		TenantID: filter.TenantID,
	})

	if err != nil {
		return false, err
	}

	if orderFound == nil {
		l.Msg("order does not exist")
		return false, nil
	}

	if !orderFound.FulfillmentStatus.CanTransitionTo(values.FulfillmentStatus) {
		l.Msg("fulfillment status cannot be transitioned to target status")
		return false, errs.InvalidTranstion(orderFound.FulfillmentStatus, values.FulfillmentStatus)
	}

	if !orderFound.PaymentStatus.CanTransitionTo(values.PaymentStatus) {
		l.Msg("payment status cannot be transitioned to target status")
		return false, errs.InvalidTranstion(orderFound.PaymentStatus, values.PaymentStatus)
	}

	if !orderFound.DeliveryStatus.CanTransitionTo(values.DeliveryStatus) {
		l.Msg("delivery status cannot be transitioned to target status")
		return false, errs.InvalidTranstion(orderFound.DeliveryStatus, values.DeliveryStatus)
	}

	l.Msg("order exists but was not patched for unknown reasons")
	return false, nil
}
