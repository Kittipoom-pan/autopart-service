package customer

import (
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(router fiber.Router, db *db.Queries, cfg *config.Config) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	// create dependencies
	repo := repository.NewCustomerRepository(db)
	usecase := usecase.NewCustomerUsecase(repo)
	controller := controller.NewCustomerController(usecase)

	router.Get("/", controller.GetAllCustomers)
	router.Get("/:id", controller.GetCustomerByID)
	router.Post("/", controller.CreateCustomer)
	router.Put("/:id", controller.UpdateCustomer)
	router.Delete("/:id", controller.DeleteCustomer)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	// repo := repository.NewCustomerRepository(db)
	// usecase := usecase.NewCustomerUsecase(repo, cfg)
	// customerController := controller.NewCustomerController(usecase)
	// authController := controller.NewAuthController(usecase) // สมมติว่ามี AuthController

	// router.Post("/login", authController.Login)
	// router.Post("/register", customerController.CreateCustomer) // การลงทะเบียนมักเป็น Public
}
