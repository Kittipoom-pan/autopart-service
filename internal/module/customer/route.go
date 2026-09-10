package customer

import (
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, jwtAuth fiber.Handler) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))
	router.Use(jwtAuth)

	repo := repository.NewCustomerRepository(db)
	customerUsecase := usecase.NewCustomerUsecase(repo)
	customerController := controller.NewCustomerController(customerUsecase)

	// Admin-only: list all customers
	router.Get("/", middleware.RequireRoles(auth.RoleSuperAdmin, auth.RoleStaff), customerController.GetAllCustomers)

	// Self or admin: read / update / delete a specific customer
	router.Get("/:id", customerController.GetCustomerByID)
	router.Put("/:id", customerController.UpdateCustomer)
	router.Delete("/:id", customerController.DeleteCustomer)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewCustomerRepository(db)
	authUsecase := usecase.NewAuthUsecase(repo, cfg)
	customerUsecase := usecase.NewCustomerUsecase(repo)
	customerController := controller.NewCustomerController(customerUsecase)
	authController := controller.NewAuthController(authUsecase, cfg)

	router.Post("/login", authController.Login)
	router.Post("/register", customerController.CreateCustomer)
}
