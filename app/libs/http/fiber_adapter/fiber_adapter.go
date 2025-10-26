package fiber_adapter

import (
	"rabi-food-core/config"
	"rabi-food-core/libs/http"
	"rabi-food-core/libs/http/controllers/tenant_controller"
	"rabi-food-core/libs/http/controllers/user_controller"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	"rabi-food-core/libs/http/routes"
	"rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	jwtware "github.com/gofiber/contrib/jwt"
)

type fiberAdapter struct {
	port string
	app  *fiber.App
}

func New(
	port string,
	tenantController *tenant_controller.TenantController,
	userController *user_controller.UserController,
) http.HTTPServer {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: middlewares.ErrorHandler,
	})

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.AuthSecret)},
	})

	requestIDMiddleware := requestid.New(requestid.Config{
		ContextKey: logger.LoggerKey,
	})

	app.
		Use(cors.New()).
		Use(requestIDMiddleware).
		Post("/tenant", tenantController.Create).
		Use(jwtMiddleware).
		Use(middlewares.Session)

	routes.User(app, userController)
	routes.Tenant(app, tenantController)

	return &fiberAdapter{
		app:  app,
		port: port,
	}
}

func (f *fiberAdapter) Start() error {
	return f.app.Listen(":" + f.port)
}

func (f *fiberAdapter) Stop() error {
	return f.app.Shutdown()
}
