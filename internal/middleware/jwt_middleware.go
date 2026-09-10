package middleware

import (
	"strings"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return helper.RespondError(c, customerror.NewUnauthorizedError("Missing or invalid authorization header"))
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.VerifyToken(tokenStr, cfg)
		if err != nil {
			return helper.RespondError(c, customerror.NewUnauthorizedError("Invalid token"))
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RequireRoles(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		role := auth.GetRole(c)
		if _, ok := allowed[role]; !ok {
			return helper.RespondError(c, customerror.NewForbiddenError("Insufficient permissions"))
		}
		return c.Next()
	}
}
