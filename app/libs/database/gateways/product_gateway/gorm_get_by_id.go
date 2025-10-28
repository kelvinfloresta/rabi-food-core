package product_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormProductGatewayAdapter) GetByID(id string) (*GetByIDOutput, error) {
	output := &models.Product{}
	result := g.DB.Conn.Limit(1).
		Preload("Category").
		Find(output, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := GetByIDOutput{
		ID:           id,
		Name:         output.Name,
		Description:  output.Description,
		Photo:        output.Photo,
		CategoryID:   output.CategoryID,
		CategoryName: output.Category.Name,
		Unit:         output.Unit,
		Price:        output.Price,
		IsActive:     output.IsActive,
	}

	return &adapted, nil
}
