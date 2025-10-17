package user_case

import g "rabi-food-core/libs/database/gateways/user_gateway"

type UserCase struct {
	gateway g.UserGateway
}

func New(gateway g.UserGateway) *UserCase {
	return &UserCase{gateway}
}
