package user_controller

import (
	"net/http"
	"rabi-food-core/app_context"
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/user_case"

	"github.com/gofiber/fiber/v2"
)

func (c *UserController) Patch(ctx *fiber.Ctx) error {
	filter := &user_case.PatchFilter{
		ID: ctx.Params("id"),
	}

	session := app_context.GetSession(ctx.Context())
	if filter.ID != session.UserID && !session.Role.IsBackoffice() {
		return ctx.SendStatus(http.StatusForbidden)
	}

	data := user_case.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	updated, err := c.usecase.Patch(ctx.Context(), *filter, data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
