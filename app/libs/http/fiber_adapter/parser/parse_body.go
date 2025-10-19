package parser

import (
	"github.com/gofiber/fiber/v2"
)

func ParseBody(c *fiber.Ctx, data any) error {
	if err := c.BodyParser(data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return err
	}

	return nil
}
