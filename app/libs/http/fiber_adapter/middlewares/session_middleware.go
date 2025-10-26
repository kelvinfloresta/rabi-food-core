package middlewares

import (
	"errors"
	"fmt"
	"rabi-food-core/app_context"
	"rabi-food-core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	errInvalidToken  = errors.New("INVALID_TOKEN")
	errInvalidClaims = errors.New("INVALID_CLAIMS")
)

// Session is a middleware that extracts user session information from the JWT token.
func Session(c *fiber.Ctx) error {
	token, ok := c.Context().UserValue("user").(*jwt.Token)
	if !ok || !token.Valid {
		return errInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errInvalidClaims
	}

	session := &app_context.UserSession{
		UserID:         fmt.Sprint(claims["user_id"]),
		TenantID:       fmt.Sprint(claims["tenant_id"]),
		Name:           fmt.Sprint(claims["name"]),
		Login:          fmt.Sprint(claims["login"]),
		OriginalUserID: fmt.Sprint(claims["original_user_id"]),
		Role:           domain.Role(fmt.Sprint(claims["role"])),
	}

	c.Context().SetUserValue(app_context.SessionKey, session)

	return c.Next()
}
