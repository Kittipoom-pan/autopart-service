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

func SetupPrivateRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, auth fiber.Handler) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	// create dependencies
	repo := repository.NewCustomerRepository(db)
	usecase := usecase.NewCustomerUsecase(repo)
	controller := controller.NewCustomerController(usecase)

	router.Get("/", auth, controller.GetAllCustomers)
	router.Get("/:id", auth, controller.GetCustomerByID)
	router.Put("/:id", auth, controller.UpdateCustomer)
	router.Delete("/:id", auth, controller.DeleteCustomer)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewCustomerRepository(db)
	authUsecase := usecase.NewAuthUsecase(repo, cfg)
	usecase := usecase.NewCustomerUsecase(repo)
	customerController := controller.NewCustomerController(usecase)
	authController := controller.NewAuthController(authUsecase, cfg)

	router.Post("/login", authController.Login)
	router.Post("/register", customerController.CreateCustomer)
}
