package tenant_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/tenant_gateway"
)

type PatchFilter struct {
	ID *string
}

type PatchValues struct {
	Name string
}

func (c TenantCase) Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error) {
	return c.gateway.Patch(
		g.PatchFilter{ID: filter.ID},
		g.PatchValues{Name: values.Name},
	)
}
