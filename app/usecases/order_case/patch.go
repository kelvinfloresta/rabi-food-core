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
	newValues g.PatchValues,
) (bool, error) {
	session := app_context.GetSession(ctx)
	if !session.Role.IsBackoffice() {
		filter.TenantID = session.TenantID
	}

	if session.Role.IsUser() {
		return false, errs.ErrForbidden
	}

	deliveryStatusIn := []order.DeliveryStatus(nil)
	if newValues.DeliveryStatus != "" {
		deliveryStatusIn = newValues.DeliveryStatus.GetPrerequisites()
		if len(deliveryStatusIn) == 0 {
			return false, errs.StatusNotModifiable(newValues.DeliveryStatus)
		}
	}

	fulfillmentStatusIn := []order.FulfillmentStatus(nil)
	if newValues.FulfillmentStatus != "" {
		fulfillmentStatusIn = newValues.FulfillmentStatus.GetPrerequisites()
		if len(fulfillmentStatusIn) == 0 {
			return false, errs.StatusNotModifiable(newValues.FulfillmentStatus)
		}
	}

	paymentStatusIn := []order.PaymentStatus(nil)
	if newValues.PaymentStatus != "" {
		paymentStatusIn = newValues.PaymentStatus.GetPrerequisites()
		if len(paymentStatusIn) == 0 {
			return false, errs.StatusNotModifiable(newValues.PaymentStatus)
		}
	}

	patched, err := c.gateway.Patch(g.PatchFilter{
		ID:                  filter.ID,
		TenantID:            filter.TenantID,
		DeliveryStatusIn:    deliveryStatusIn,
		FulfillmentStatusIn: fulfillmentStatusIn,
		PaymentStatusIn:     paymentStatusIn,
	}, newValues)

	if err != nil {
		return false, err
	}

	if patched {
		return true, nil
	}

	return c.handleNotPatched(ctx, filter, newValues)
}

func (c *OrderCase) handleNotPatched(ctx context.Context, filter PatchFilter, newValues g.PatchValues) (bool, error) {
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

	if !orderFound.FulfillmentStatus.CanTransitionTo(newValues.FulfillmentStatus) {
		l.Msg("fulfillment status cannot be transitioned to target status")
		return false, errs.InvalidTranstion(orderFound.FulfillmentStatus, newValues.FulfillmentStatus)
	}

	if !orderFound.PaymentStatus.CanTransitionTo(newValues.PaymentStatus) {
		l.Msg("payment status cannot be transitioned to target status")
		return false, errs.InvalidTranstion(orderFound.PaymentStatus, newValues.PaymentStatus)
	}

	if !orderFound.DeliveryStatus.CanTransitionTo(newValues.DeliveryStatus) {
		l.Msg("delivery status cannot be transitioned to target status")
		return false, errs.InvalidTranstion(orderFound.DeliveryStatus, newValues.DeliveryStatus)
	}

	l.Msg("order exists but was not patched for unknown reasons")
	return false, nil
}
