package factories

import (
	"rabi-food-core/libs/database"
	g "rabi-food-core/libs/database/gateways/tenant_gateway"
	"rabi-food-core/libs/database/gorm_adapter"
	c "rabi-food-core/libs/http/controllers/tenant_controller"

	"rabi-food-core/usecases/tenant_case"
)

func NewTenant(d database.Database) *c.TenantController {
	DB, ok := d.(*gorm_adapter.GormAdapter)
	if !ok {
		panic(ErrDatabaseAdapter)
	}

	gateway := &g.GormTenantGatewayAdapter{DB: DB}
	userCase := newUserCase(DB)

	usecase := tenant_case.New(gateway, userCase)
	return c.New(usecase)
}
