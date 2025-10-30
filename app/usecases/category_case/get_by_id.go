package category_case

import (
	"context"
	"rabi-food-core/app_context"
	g "rabi-food-core/libs/database/gateways/category_gateway"
)

func (c *CategoryCase) GetByID(ctx context.Context, id string) (*g.GetByIDOutput, error) {
	session := app_context.GetSession(ctx)
	return c.gateway.GetByID(id, session.TenantID)
}
