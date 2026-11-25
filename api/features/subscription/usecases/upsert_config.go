package usecases

import (
	"context"

	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
)

func (s *SubscriptionCase) UpsertConfig(ctx context.Context, input g.UpsertConfigInput) error {
	return s.gateway.UpsertConfig(ctx, input.TenantID, input)
}
