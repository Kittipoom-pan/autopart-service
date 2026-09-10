package utils

import (
	apperror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func GetUserID(c *fiber.Ctx, logger zerolog.Logger) (int, error) {
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		logger.Error().Msg("user_id not found in fiber locals")
		return 0, apperror.NewUnauthorizedError("Invalid token")
	}

	return int(userIDRaw.(uint32)), nil
}
