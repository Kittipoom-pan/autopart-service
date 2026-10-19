package server

import (
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")

	customerPublic := v1.Group("/customer")
	customer.SetupPublicRoutes(customerPublic, s.Db, s.Cfg, s.Logger)

	customerPrivate := v1.Group("/customer")
	customer.SetupPrivateRoutes(
		customerPrivate, s.Db, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	adminPublic := v1.Group("/admin")
	admin.SetupPublicRoutes(adminPublic, s.Db, s.Cfg, s.Logger)

	adminPrivate := v1.Group("/admin")
	admin.SetupPrivateRoutes(
		adminPrivate, s.Db, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	partPublic := v1.Group("/part")
	part.SetupPublicRoutes(partPublic, s.Db, s.Cfg, s.Logger)

	partPrivate := v1.Group("/part")
	part.SetupPrivateRoutes(
		partPrivate, s.Db, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	orderPublic := v1.Group("/order")
	order.SetupPublicRoutes(
		orderPublic, s.Conn, s.Cfg, s.Logger,
	)

	orderPrivate := v1.Group("/order")
	order.SetupPrivateRoutes(
		orderPrivate, s.Conn, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Endpoint not found",
		})
	})

	return nil
}
