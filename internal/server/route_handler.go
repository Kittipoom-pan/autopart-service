package server

import (
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")

	// public routes (no authentication required)
	customerPublic := v1.Group("/customer")
	customer.SetupPublicRoutes(customerPublic, s.Db, s.Cfg)
	// private routes (protected by JWT authentication)
	customerPrivate := v1.Group("/customer")
	customerPrivate.Use(middleware.JWTMiddleware(s.Cfg))
	customer.SetupPrivateRoutes(customerPrivate, s.Db, s.Cfg)

	adminPublic := v1.Group("/admin")
	admin.SetupPublicRoutes(adminPublic, s.Db, s.Cfg)

	adminPrivate := v1.Group("/admin")
	adminPrivate.Use(middleware.JWTMiddleware(s.Cfg))
	admin.SetupPrivateRoutes(adminPrivate, s.Db, s.Cfg)

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
