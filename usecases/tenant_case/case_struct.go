package tenant_case

import (
	g "rabi-food-core/libs/database/gateways/tenant_gateway"
	"rabi-food-core/usecases/user_case"
)

type TenantCase struct {
	gateway  g.TenantGateway
	userCase *user_case.UserCase
}

func New(
	gateway g.TenantGateway,
	userCase *user_case.UserCase,
) TenantCase {
	return TenantCase{
		gateway:  gateway,
		userCase: userCase,
	}
}
