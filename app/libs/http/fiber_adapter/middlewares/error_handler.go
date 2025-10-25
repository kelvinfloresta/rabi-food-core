package middlewares

import (
	"rabi-food-core/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(struct {
			Errors []string `json:"errors"`
		}{
			Errors: utils.TranslateValidationErrors(errs),
		})
	}

	return err
}
