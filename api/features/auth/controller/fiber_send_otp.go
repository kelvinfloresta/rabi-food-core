package controller

import (
	"net/http"

	"github.com/Rabi-IT/rabi-food-core/features/auth/usecases"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"
	"github.com/Rabi-IT/rabi-food-core/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func (c *AuthController) SendOTP(ctx *fiber.Ctx) error {
	data := &usecases.SendOTPInput{}
	err := ctx.BodyParser(data)
	if err != nil {
		return ctx.JSON(err)
	}

	err = validator.V.Struct(data)
	if err != nil {
		return err
	}

	uctx := ctx.UserContext()
	logger.GetWideEvent(uctx).Event = "auth-send-otp"

	err = c.usecase.SendOTP(uctx, data)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
