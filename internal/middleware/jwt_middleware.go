package middleware

import (
	"strings"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			log.Warn().
				Str("path", c.Path()).
				Msg("Missing Authorization Header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization"})
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.VerifyToken(tokenStr, cfg)
		if err != nil {
			log.Error().
				Err(err).
				Str("path", c.Path()).
				Str("token", tokenStr).
				Msg("JWT verification failed")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// add claims to context Fiber
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}
