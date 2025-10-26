package middlewares

import (
	"rabi-food-core/libs/logger"
	lib "rabi-food-core/libs/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationErrorResponse struct {
	Errors []lib.ValidationError `json:"errors"`
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
			Errors: lib.ParseValidationError(errs),
		})
	}

	logger.L().Error().Err(err).Msg("internal server error")

	return err
}
