package product_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/product_gateway"
)

func (c *ProductCase) Patch(
	ctx context.Context,
	filter g.PatchFilter,
	values g.PatchValues,
) (bool, error) {
	return c.gateway.Patch(filter, values)
}
