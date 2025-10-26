package middlewares

import (
	"rabi-food-core/libs/logger"
	lib "rabi-food-core/libs/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ValidationErrorResponse represents the structure of validation error responses.
// It is used to return detailed information about validation errors to the client.
type ValidationErrorResponse struct {
	Errors []lib.ValidationError `json:"errors"`
}

// ErrorHandler is a middleware that handles errors occurring during request processing.
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	//nolint:errorlint
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(ValidationErrorResponse{
			Errors: lib.ParseValidationError(errs),
		})
	}

	logger.L().Error().Err(err).Msg("internal server error")

	return err
}
