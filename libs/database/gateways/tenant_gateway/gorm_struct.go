package tenant_gateway

import "rabi-food-core/libs/database/gorm_adapter"

type GormTenantGatewayAdapter struct {
	DB *gorm_adapter.GormAdapter
}
