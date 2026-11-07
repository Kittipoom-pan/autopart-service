package customer

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func SetupPrivateRoutes(router fiber.Router, conn *sql.DB, cfg *config.Config, jwtAuth fiber.Handler, logger zerolog.Logger) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))
	router.Use(jwtAuth)

	txManager := database.NewTxManager(conn, logger)
	repo := repository.NewCustomerRepository(conn, logger)
	customerUsecase := usecase.NewCustomerUsecase(repo, txManager, logger)
	customerController := controller.NewCustomerController(customerUsecase, logger)

	// Admin-only: list all customers
	router.Get("/", middleware.RequireRoles(auth.RoleSuperAdmin, auth.RoleStaff), customerController.GetAllCustomers)

	// Self or admin: read / update / delete a specific customer
	router.Get("/:id", customerController.GetCustomerByID)
	router.Put("/:id", customerController.UpdateCustomer)
	router.Delete("/:id", customerController.DeleteCustomer)
}

func SetupPublicRoutes(router fiber.Router, conn *sql.DB, cfg *config.Config, logger zerolog.Logger) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	txManager := database.NewTxManager(conn, logger)
	repo := repository.NewCustomerRepository(conn, logger)
	authUsecase := usecase.NewAuthUsecase(repo, cfg, logger)
	customerUsecase := usecase.NewCustomerUsecase(repo, txManager, logger)
	customerController := controller.NewCustomerController(customerUsecase, logger)
	authController := controller.NewAuthController(authUsecase, cfg, logger)

	router.Post("/login", authController.Login)
	router.Post("/register", customerController.CreateCustomer)
}
