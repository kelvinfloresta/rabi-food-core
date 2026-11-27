package controller

import (
	"fmt"
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/app_context"
	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"
	"github.com/gofiber/fiber/v2"
)

type createInput struct {
	Name string `json:"name" validate:"required"`
}

func (c *TenantController) Create(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "create-tenant"

	data := createInput{}
	err := ctx.BodyParser(&data)
	if err != nil {
		return fmt.Errorf("%w: %w", errs.ErrInvalidBody, err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	session := app_context.GetSession(uctx)
	id, err := c.usecase.Create(uctx, data.Name, session.UserID)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).SendString(id)
}
