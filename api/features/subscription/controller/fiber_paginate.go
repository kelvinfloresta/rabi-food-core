package controller

import (
	"github.com/Rabi-IT/rabi-food-core/app_context"
	g "github.com/Rabi-IT/rabi-food-core/features/subscription/gateway"
	"github.com/Rabi-IT/rabi-food-core/libs/database"

	"github.com/gofiber/fiber/v2"
)

// Paginate godoc
// @Summary Paginate subscriptions
// @Description Paginate subscriptions using query filters
// @Tags subscriptions
// @Produce json
// @Param Page query int false "Page number"
// @Param PageSize query int false "Page size"
// @Success 200 {object} g.PaginateOutput
// @Failure 500 {string} string "Internal server error"
// @Router /subscription/ [get].
func (c *SubscriptionController) Paginate(ctx *fiber.Ctx) error {
	uctx := ctx.UserContext()
	session := app_context.GetSession(uctx)

	// Tenant staff see every subscription of their tenant; anyone else only their own, since
	// customers may use any tenant through the X-Tenant-ID header.
	filter := g.PaginateFilter{TenantID: session.TenantID}
	if !session.IsTenant() {
		filter.TenantID = ctx.Get("X-Tenant-ID")
		filter.UserID = session.UserID
	}

	paginate := database.PaginateInput{
		Page:     ctx.QueryInt("Page", database.DefaultPage),
		PageSize: ctx.QueryInt("PageSize", database.DefaultPageSize),
	}

	result, err := c.usecase.Paginate(uctx, filter, paginate)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
