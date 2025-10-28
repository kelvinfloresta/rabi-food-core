package category_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/category_gateway"
)

func (c *CategoryCase) Patch(
	ctx context.Context,
	filter g.PatchFilter,
	values g.PatchValues,
) (bool, error) {
	return c.gateway.Patch(filter, values)
}
