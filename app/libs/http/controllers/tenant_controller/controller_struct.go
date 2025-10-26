package tenant_controller

import (
	"rabi-food-core/usecases/tenant_case"
)

type TenantController struct {
	usecase *tenant_case.TenantCase
}

func New(usecase *tenant_case.TenantCase) *TenantController {
	return &TenantController{usecase}
}
