package product_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/product_gateway"
)

func (c *ProductCase) List(ctx context.Context, ids []string) ([]g.ListOutput, error) {
	filter := g.ListFilter{
		IDs: ids,
	}

	session := app_context.GetSession(ctx)
	if session.Role.IsUser() {
		filter.TenantID = session.TenantID
	}

	return c.gateway.List(filter)
}
