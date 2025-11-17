package order_controller

import (
	"net/http"
	"rabi-food-core/libs/database/gateways/order_gateway"
	"rabi-food-core/libs/http/fiber_adapter/parser"
	"rabi-food-core/libs/validator"
	"rabi-food-core/usecases/order_case"

	"github.com/gofiber/fiber/v2"
)

func (c *OrderController) Patch(ctx *fiber.Ctx) error {
	filter := order_case.PatchFilter{
		ID: ctx.Params("id"),
	}

	data := order_gateway.PatchValues{}
	err := parser.ParseBody(ctx, &data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	updated, err := c.usecase.Patch(ctx.Context(), filter, data)

	if err != nil {
		return err
	}

	if updated {
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusNotFound)
}
