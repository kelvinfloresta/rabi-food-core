package tenant_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/domain"
	g "rabi-food-core/libs/database/gateways/tenant_gateway"
	"rabi-food-core/usecases/user_case"
)

type CreateInput struct {
	Name     string `validate:"required"`
	UserName string `validate:"required"`
	Phone    string `validate:"required"`
	Email    string `validate:"required,email"`
}

type CreateOutput struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func (c *TenantCase) Create(ctx context.Context, input CreateInput) (out CreateOutput, err error) {
	tenantId, err := c.gateway.Create(g.CreateInput{
		Name: input.Name,
	})

	if err != nil {
		return
	}

	ctx = app_context.WithSession(ctx, &app_context.UserSession{
		TenantID: tenantId,
		Role:     domain.UserRole,
	})
	userId, err := c.userCase.Create(ctx, &user_case.CreateInput{
		Name:  input.UserName,
		Phone: input.Phone,
		Email: input.Email,
	})

	if err != nil {
		return
	}

	return CreateOutput{
		UserID: userId,
		ID:     tenantId,
	}, nil
}
