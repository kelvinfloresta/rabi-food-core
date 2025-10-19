package routes

import (
	"rabi-food-core/libs/http/controllers/user_controller"

	"github.com/gofiber/fiber/v2"
)

func User(app *fiber.App, c *user_controller.UserController) {
	route := app.Group("/user")
	route.Post("/", c.Create)
	route.Patch("/:id", c.Patch)
	route.Delete("/:id", c.Delete)
	route.Get("/:id", c.GetByID)
	route.Get("/", c.Paginate)
}
