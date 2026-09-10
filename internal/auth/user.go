package auth

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

	userID, ok := userIDRaw.(uint32)
	if !ok {
		logger.Error().Interface("user_id", userIDRaw).Msg("user_id has unexpected type")
		return 0, apperror.NewUnauthorizedError("Invalid token")
	}

	return int(userID), nil
}

func GetRole(c *fiber.Ctx) string {
	roleRaw := c.Locals("role")
	if roleRaw == nil {
		return ""
	}
	role, _ := roleRaw.(string)
	return role
}

func EnsureSelfOrAdmin(c *fiber.Ctx, resourceID int, logger zerolog.Logger) error {
	userID, err := GetUserID(c, logger)
	if err != nil {
		return err
	}

	if IsAdminRole(GetRole(c)) {
		return nil
	}

	if userID != resourceID {
		return apperror.NewForbiddenError("You can only access your own resource")
	}

	return nil
}
