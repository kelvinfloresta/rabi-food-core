package category_case

import (
	"context"
	"rabi-food-core/libs/database"
	g "rabi-food-core/libs/database/gateways/category_gateway"
)

func (c *CategoryCase) Paginate(
	ctx context.Context,
	input g.PaginateFilter,
	paginate database.PaginateInput,
) (g.PaginateOutput, error) {
	return c.gateway.Paginate(input, paginate)
}
