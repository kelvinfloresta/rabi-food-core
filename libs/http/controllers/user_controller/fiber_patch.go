package user_controller

import (
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/usecases/user_case"

	"github.com/gofiber/fiber/v2"
)

func (c *UserController) Patch(ctx *fiber.Ctx) error {
	filter := &user_case.PatchFilter{
		ID: ctx.Params("id"),
	}

	data := user_case.PatchValues{}
	if err := parser.ParseBody(ctx, &data); err != nil {
		return ctx.JSON(err)
	}

	updated, err := c.usecase.Patch(ctx.Context(), *filter, data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(200)
	}

	return ctx.SendStatus(404)
}
