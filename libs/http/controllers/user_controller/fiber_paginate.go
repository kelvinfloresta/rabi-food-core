package user_controller

import (
	"rabi-food-core/libs/database"
	"rabi-food-core/usecases/user_case"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (c *UserController) Paginate(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("Page", "0"))
	if err != nil {
		return err
	}

	pageSize, err := strconv.Atoi(ctx.Query("PageSize", "10"))
	if err != nil {
		return err
	}

	filter := user_case.PaginateFilter{}
	if err = ctx.QueryParser(&filter); err != nil {
		return err
	}

	paginate := database.PaginateInput{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.usecase.Paginate(ctx.Context(), filter, paginate)

	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
