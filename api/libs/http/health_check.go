package http

import (
	"context"
	"time"

	auth_gateway "github.com/Rabi-IT/rabi-food-core/features/auth/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
	"github.com/Rabi-IT/rabi-food-core/libs/logger"

	"github.com/gofiber/fiber/v2"
)

// healthCheckTimeout is applied per dependency, keeping the whole check under the
// ALB health check timeout (5s by default).
const healthCheckTimeout = 2 * time.Second

// HealthCheck verifies the services the API cannot operate without.
type HealthCheck struct {
	db   *database.PgxAdapter
	auth auth_gateway.AuthGateway
}

func NewHealthCheck(db *database.PgxAdapter, auth auth_gateway.AuthGateway) *HealthCheck {
	return &HealthCheck{
		db:   db,
		auth: auth,
	}
}

// Handle answers with a bare 200 when the database and GoTrue respond, or 503 naming the
// unavailable ones. The underlying errors are logged, never exposed in the response.
func (h *HealthCheck) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	unavailable := []string{}

	if err := probe(ctx, h.db.Ping); err != nil {
		logger.Get(ctx).Error().Err(err).Msg("health check failed: database")
		unavailable = append(unavailable, "database")
	}

	if err := probe(ctx, h.auth.Health); err != nil {
		logger.Get(ctx).Error().Err(err).Msg("health check failed: auth")
		unavailable = append(unavailable, "auth")
	}

	if len(unavailable) > 0 {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"unavailable": unavailable})
	}

	return c.SendStatus(fiber.StatusOK)
}

func probe(ctx context.Context, check func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()

	return check(ctx)
}
