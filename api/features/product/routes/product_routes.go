package routes

import (
	"github.com/Rabi-IT/rabi-food-core/features/product/controller"
	"github.com/Rabi-IT/rabi-food-core/libs/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

// Product registers public product routes (before JWT middleware).
func Product(app *fiber.App, c *controller.ProductController) {
	route := app.Group("/product")
	route.Get("/catalog", middlewares.RequireSelectedTenant, c.Catalog)
	route.Get("/:id", c.GetByID)
}

// ProductProtected registers product routes that require a valid JWT.
func ProductProtected(app *fiber.App, c *controller.ProductController) {
	route := app.Group("/product")
	route.Post("/", middlewares.RequireTenantMember, c.Create)
	route.Patch("/:id", middlewares.RequireTenantMember, c.Patch)
	route.Delete("/:id", middlewares.RequireTenantMember, c.Delete)

	backoffice := app.Group("/backoffice/product", middlewares.RequireBackoffice)
	backoffice.Get("/", c.BackofficePaginate)
}
