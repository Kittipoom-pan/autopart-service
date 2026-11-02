package server

import (
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart"
	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method"
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
		orderPublic, s.Conn, s.Db, s.Cfg, s.Logger,
	)

	orderPrivate := v1.Group("/order")
	order.SetupPrivateRoutes(
		orderPrivate, s.Conn, s.Db, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	cartPrivate := v1.Group("/cart")
	cart.SetupPrivateRoutes(
		cartPrivate, s.Conn, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	cartRoot := s.App.Group("/cart")
	cart.SetupPrivateRoutes(
		cartRoot, s.Conn, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	addressPrivate := v1.Group("/address")
	address.SetupPrivateRoutes(
		addressPrivate, s.Conn, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	paymentMethodPrivate := v1.Group("/payment-method")
	customer_payment_method.SetupPrivateRoutes(
		paymentMethodPrivate, s.Conn, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	checkoutPrivate := v1.Group("/checkout")
	checkout.SetupPrivateRoutes(
		checkoutPrivate, s.Conn, s.Db, s.Cfg,
		middleware.JWTMiddleware(s.Cfg), s.Logger,
	)

	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Endpoint not found",
		})
	})

	return nil
}
