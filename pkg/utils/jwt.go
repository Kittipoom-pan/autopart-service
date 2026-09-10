package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromToken(c *fiber.Ctx) (uint32, error) {
	// get user from Locals (default value Fiber JWT is user)
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return 0, errors.New("invalid or missing token")
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("user id not found in token")
	}

	return uint32(idFloat), nil
}
