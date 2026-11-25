package middlewares

import (
	"github.com/Rabi-IT/rabi-food-core/libs/errs"

	"github.com/gofiber/fiber/v2"
)

// RequireSelectedTenant requires the tenant a customer selected (X-Tenant-ID header).
func RequireSelectedTenant(c *fiber.Ctx) error {
	if c.Get("X-Tenant-ID") == "" {
		return errs.ErrMissingTenantID
	}

	return c.Next()
}
