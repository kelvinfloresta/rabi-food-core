package middlewares

import (
	"rabi-food-core/libs/logger"
	"rabi-food-core/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationErrorResponse struct {
	Errors []utils.ValidationError `json:"errors"`
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
			Errors: utils.ParseValidationError(errs),
		})
	}

	logger.L().Error().Err(err).Msg("internal server error")

	return err
}
