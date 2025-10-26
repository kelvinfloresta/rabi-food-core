package parser

import (
	"github.com/gofiber/fiber/v2"
)

func ParseBody(c *fiber.Ctx, data any) error {
	err := c.BodyParser(data)

	if err != nil {
		c.Status(fiber.StatusBadRequest)

		return err
	}

	return nil
}
