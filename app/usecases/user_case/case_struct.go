package user_case

import g "rabi-food-core/libs/database/gateways/user_gateway"

// UserCase encapsulates the business logic related to users.
type UserCase struct {
	gateway g.UserGateway
}

// New creates a new instance of UserCase with the provided UserGateway.
func New(gateway g.UserGateway) *UserCase {
	return &UserCase{gateway}
}
