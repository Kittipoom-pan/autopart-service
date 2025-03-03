package user

import (
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/module/user/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/user/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/user/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router, db *sql.DB) {
	// Create dependencies
	repo := repository.NewUserRepository(db)
	usecase := usecase.NewUsersUsecase(repo)
	controller := controller.NewUsersController(usecase)

	router.Get("/", controller.GetUsers)
	router.Post("/", controller.CreateUser)
}
