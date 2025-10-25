package tenant_controller

import (
	"rabi-food-core/usecases/tenant_case"
	"rabi-food-core/utils"

	"github.com/gofiber/fiber/v2"
)

func (c *TenantController) Create(ctx *fiber.Ctx) error {
	data := tenant_case.CreateInput{}
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.JSON(err)
	}

	if err := utils.Validator.Struct(data); err != nil {
		return err
	}

	output, err := c.usecase.Create(ctx.Context(), data)

	if err != nil {
		return err
	}

	return ctx.Status(201).JSON(output)
}
