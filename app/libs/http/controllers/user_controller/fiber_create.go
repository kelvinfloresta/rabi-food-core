package user_controller

import (
	"rabi-food-core/usecases/user_case"
	"rabi-food-core/utils"

	"github.com/gofiber/fiber/v2"
)

func (c *UserController) Create(ctx *fiber.Ctx) error {
	data := &user_case.CreateInput{}
	if err := ctx.BodyParser(data); err != nil {
		return ctx.JSON(err)
	}

	if err := utils.Validator.Struct(data); err != nil {
		return err
	}

	id, err := c.usecase.Create(ctx.Context(), data)

	if err != nil {
		return err
	}

	return ctx.Status(201).SendString(id)
}
