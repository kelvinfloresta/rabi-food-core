package tenant_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormTenantGatewayAdapter) Patch(filter PatchFilter, newValues PatchValues) (bool, error) {
	query := g.DB.Conn.Model(&models.Tenant{}).Where("id = ?", filter.ID)

	result := query.Updates(newValues)

	return result.RowsAffected > 0, result.Error
}
