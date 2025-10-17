package factories

import (
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/http/controllers/user_controller"

	"rabi-food-core/usecases/user_case"
)

func newUserCase(d database.Database) *user_case.UserCase {
	DB, ok := d.(*gorm_adapter.GormAdapter)
	if !ok {
		panic(ErrDatabaseAdapter)
	}

	gateway := &user_gateway.GormUserGatewayAdapter{DB: DB}
	return user_case.New(gateway)
}

func NewUser(d database.Database) *user_controller.UserController {
	DB, ok := d.(*gorm_adapter.GormAdapter)
	if !ok {
		panic(ErrDatabaseAdapter)
	}

	usecase := newUserCase(DB)
	return user_controller.New(usecase)
}
