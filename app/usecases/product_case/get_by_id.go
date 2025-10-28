package product_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/product_gateway"
)

func (c *ProductCase) GetByID(ctx context.Context, id string) (*g.GetByIDOutput, error) {
	return c.gateway.GetByID(id)
}
