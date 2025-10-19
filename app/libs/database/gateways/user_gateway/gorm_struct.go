package user_gateway

import "rabi-food-core/libs/database/gorm_adapter"

type GormUserGatewayAdapter struct {
	DB *gorm_adapter.GormAdapter
}
