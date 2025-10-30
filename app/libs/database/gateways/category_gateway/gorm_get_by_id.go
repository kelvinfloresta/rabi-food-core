package category_gateway

import (
	"rabi-food-core/libs/database/gorm_adapter/models"
)

func (g *GormCategoryGatewayAdapter) GetByID(id string, tenantID string) (*GetByIDOutput, error) {
	output := &models.Category{}
	result := g.DB.Conn.Limit(1).Find(output, "id = ? AND tenant_id = ?", id, tenantID)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	adapted := GetByIDOutput{
		ID:          output.ID,
		Name:        output.Name,
		Description: output.Description,
		TenantID:    output.TenantID,
	}

	return &adapted, nil
}
