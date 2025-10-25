package user_case_test

import (
	"rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/usecases/user_case"
)

func makeSut(g user_gateway.UserGateway) *user_case.UserCase {
	return user_case.New(g)
}
