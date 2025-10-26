package di

import (
	"errors"
	"rabi-food-core/config"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/database/gateways/tenant_gateway"
	"rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/http"
	"rabi-food-core/libs/http/controllers/tenant_controller"
	"rabi-food-core/libs/http/controllers/user_controller"
	"rabi-food-core/libs/http/fiber_adapter"
	"rabi-food-core/usecases/tenant_case"
	"rabi-food-core/usecases/user_case"

	"github.com/samber/do"
)

func newInjector(dbConfig *config.DatabaseConfig) *do.Injector {
	injector := do.New()

	// Database
	do.Provide(injector, func(i *do.Injector) (*gorm_adapter.GormAdapter, error) {
		db, ok := gorm_adapter.New(dbConfig).(*gorm_adapter.GormAdapter)
		if !ok {
			return nil, errors.New("failed to cast database adapter")
		}
		return db, nil
	})

	// Database as interface
	do.Provide(injector, func(i *do.Injector) (database.Database, error) {
		return do.MustInvoke[*gorm_adapter.GormAdapter](i), nil
	})

	// HTTP Server
	do.Provide(injector, func(i *do.Injector) (http.HTTPServer, error) {
		tenantController := do.MustInvoke[*tenant_controller.TenantController](i)
		userController := do.MustInvoke[*user_controller.UserController](i)
		return fiber_adapter.New(config.Port, tenantController, userController), nil
	})

	// User dependencies
	do.Provide(injector, func(i *do.Injector) (user_gateway.UserGateway, error) {
		db := do.MustInvoke[*gorm_adapter.GormAdapter](i)
		return &user_gateway.GormUserGatewayAdapter{DB: db}, nil
	})

	do.Provide(injector, func(i *do.Injector) (*user_case.UserCase, error) {
		gw := do.MustInvoke[user_gateway.UserGateway](i)
		return user_case.New(gw), nil
	})

	do.Provide(injector, func(i *do.Injector) (*user_controller.UserController, error) {
		c := do.MustInvoke[*user_case.UserCase](i)
		return user_controller.New(c), nil
	})

	// Tenant dependencies
	do.Provide(injector, func(i *do.Injector) (tenant_gateway.TenantGateway, error) {
		db := do.MustInvoke[*gorm_adapter.GormAdapter](i)
		return &tenant_gateway.GormTenantGatewayAdapter{DB: db}, nil
	})

	do.Provide(injector, func(i *do.Injector) (*tenant_case.TenantCase, error) {
		gw := do.MustInvoke[tenant_gateway.TenantGateway](i)
		userCase := do.MustInvoke[*user_case.UserCase](i)
		return tenant_case.New(gw, userCase), nil
	})

	do.Provide(injector, func(i *do.Injector) (*tenant_controller.TenantController, error) {
		c := do.MustInvoke[*tenant_case.TenantCase](i)
		return tenant_controller.New(c), nil
	})

	return injector
}

func NewProduction() *do.Injector {
	return newInjector(config.ProductionDatabase)
}

func NewTest() *do.Injector {
	return newInjector(config.TestDatabase)
}
