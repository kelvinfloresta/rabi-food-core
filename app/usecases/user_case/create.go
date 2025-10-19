package user_case

import (
	"context"
	"rabi-food-core/app_context"
	"rabi-food-core/domain"
	g "rabi-food-core/libs/database/gateways/user_gateway"
	"rabi-food-core/utils"
)

type CreateInput struct {
	Name         string `validate:"required"`
	Email        string `validate:"required,email"`
	TenantID     string
	Photo        string
	TaxID        string
	City         string
	State        string
	Phone        string
	ZIP          string
	SocialID     string
	Neighborhood string
	Street       string
	Complement   string
}

func (c *UserCase) Create(ctx context.Context, input *CreateInput) (string, error) {
	if err := utils.Validator.Struct(input); err != nil {
		return "", err
	}
	session := app_context.GetSession(ctx)

	tenantId := ""
	if session.Role.IsUser() {
		tenantId = session.TenantID
	} else if session.Role.IsBackoffice() {
		tenantId = input.TenantID
	}

	return c.gateway.Create(g.CreateInput{
		TenantID:     tenantId,
		City:         input.City,
		State:        input.State,
		ZIP:          input.ZIP,
		Phone:        input.Phone,
		Email:        input.Email,
		Photo:        input.Photo,
		TaxID:        input.TaxID,
		SocialID:     input.SocialID,
		Street:       input.Street,
		Complement:   input.Complement,
		Neighborhood: input.Neighborhood,
		Name:         input.Name,
		Role:         domain.UserRole,
	})
}
