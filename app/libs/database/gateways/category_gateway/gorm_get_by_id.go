package category_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormCategoryGatewayAdapter) GetByID(id string) (*GetByIDOutput, error) {
	output := &models.Category{}
	result := g.DB.Conn.Limit(1).Find(output, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := GetByIDOutput{
		ID:          id,
		Name:        output.Name,
		Description: output.Description,
	}

	return &adapted, nil
}
