package tenant_case

import (
	"context"
	g "rabi-food-core/libs/database/gateways/tenant_gateway"
)

func (c *TenantCase) GetByID(ctx context.Context, id string) (*g.GetByIDOutput, error) {
	return c.gateway.GetByID(id)
}
