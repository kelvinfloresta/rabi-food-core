package controller

import (
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/auth/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) SwitchTenant(ctx *fiber.Ctx) error {
	data := &usecases.SwitchTenantInput{}
	if err := ctx.BodyParser(data); err != nil {
		return fmt.Errorf("%w: %w", errs.ErrInvalidBody, err)
	}

	if err := validator.V.Struct(data); err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "auth-switch-tenant"

	out, err := c.usecase.SwitchTenant(uctx, data)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(out)
}
