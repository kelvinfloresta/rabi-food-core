package product_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormProductGatewayAdapter) Delete(id string) (bool, error) {
	result := g.DB.Conn.Where(
		"id = ?", id,
	).Delete(&models.Product{})

	return result.RowsAffected > 0, result.Error
}
