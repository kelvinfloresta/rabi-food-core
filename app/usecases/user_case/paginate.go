package user_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/libs/database"
	g "rabi-food-core/libs/database/gateways/user_gateway"
)

type PaginateFilter struct {
	State *string
	City  *string
	Name  *string
}

var EMPTY_PAGINATION = g.PaginateOutput{
	Data: []g.PaginateData{},
}

func (c *UserCase) Paginate(ctx context.Context, input PaginateFilter, paginate database.PaginateInput) (g.PaginateOutput, error) {
	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		return EMPTY_PAGINATION, nil
	}

	return c.gateway.Paginate(g.PaginateFilter{
		City:  input.City,
		State: input.State,
		Name:  input.Name,
	}, paginate)
}
