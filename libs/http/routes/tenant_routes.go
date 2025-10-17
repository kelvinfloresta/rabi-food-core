package routes

import (
	"rabi-food-core/libs/http/controllers/tenant_controller"

	"github.com/gofiber/fiber/v2"
)

func Tenant(app *fiber.App, c *tenant_controller.TenantController) {
	route := app.Group("/tenant")
	route.Get("/:id", c.GetByID)
	route.Patch("/:id", c.Patch)
}
