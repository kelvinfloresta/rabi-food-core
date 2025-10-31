package category_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormCategoryGatewayAdapter) Delete(filter DeleteFilter) (bool, error) {
	result := g.DB.Conn.Where(
		"id = ? AND tenant_id = ?", filter.ID, filter.TenantID,
	).Delete(&models.Category{})

	return result.RowsAffected > 0, result.Error
}
