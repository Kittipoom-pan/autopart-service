package customer

import (
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router, db *db.Queries) {
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
