package product_case

import (
	"context"
	"rabi-food-core/libs/database"
	g "rabi-food-core/libs/database/gateways/product_gateway"
)

func (c *ProductCase) Paginate(
	ctx context.Context,
	input g.PaginateFilter,
	paginate database.PaginateInput,
) (g.PaginateOutput, error) {
	return c.gateway.Paginate(input, paginate)
}
