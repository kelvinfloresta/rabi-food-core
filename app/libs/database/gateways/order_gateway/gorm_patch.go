package order_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormOrderGatewayAdapter) Patch(filter PatchFilter, newValues PatchValues) (bool, error) {
	query := g.DB.Conn.Model(&models.Order{}).Where("id = ?", filter.ID)

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	if len(filter.DeliveryStatusIn) > 0 {
		query = query.Where("delivery_status IN ?", filter.DeliveryStatusIn)
	}

	if len(filter.FulfillmentStatusIn) > 0 {
		query = query.Where("fulfillment_status IN ?", filter.FulfillmentStatusIn)
	}

	if len(filter.PaymentStatusIn) > 0 {
		query = query.Where("payment_status IN ?", filter.PaymentStatusIn)
	}

	result := query.Updates(newValues)

	return result.RowsAffected > 0, result.Error
}
