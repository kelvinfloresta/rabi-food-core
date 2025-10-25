package fiber_adapter

import (
	"rabi-food-core/config"
	"rabi-food-core/factories"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/http"
	"rabi-food-core/libs/http/controllers/auth_controller"
	"rabi-food-core/libs/http/fiber_adapter/middlewares"
	"rabi-food-core/libs/http/routes"
	"rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	jwtware "github.com/gofiber/contrib/jwt"
)

type fiberAdapter struct {
	app *fiber.App
}

func New(d database.Database) http.HTTPServer {
	return newFiber(d)
}

func (f *fiberAdapter) Start(port string) error {
	return f.app.Listen(":" + port)
}

func (f *fiberAdapter) Stop() error {
	return f.app.Shutdown()
}

func newFiber(d database.Database) http.HTTPServer {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: middlewares.ErrorHandler,
	})

	tenantController := factories.NewTenant(d)
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
		Use(auth_controller.Session)

	routes.User(app, factories.NewUser(d))
	routes.Tenant(app, tenantController)

	return &fiberAdapter{app}
}
