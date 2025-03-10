package server

import (
	user "github.com/Kittipoom-pan/autopart-service/internal/module/customer"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")

	usersGroup := v1.Group("/customer")
	user.SetupRoutes(usersGroup, s.Db)

	// Products routes
	// ...

	// End point not found
	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     "error, end point not found",
			"result":      nil,
		})
	})

	return nil
}
