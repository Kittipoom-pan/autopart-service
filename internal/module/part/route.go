package part

import (
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, auth fiber.Handler) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewPartRepository(db)
	usecase := usecase.NewPartUsecase(repo)
	controller := controller.NewPartController(usecase)

	router.Put("/:id", auth, controller.UpdatePart)
	router.Post("/", controller.CreatePart)
	router.Delete("/:id", auth, controller.DeletePart)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewPartRepository(db)
	usecase := usecase.NewPartUsecase(repo)
	partController := controller.NewPartController(usecase)

	router.Get("/", partController.GetAllParts)
	router.Get("/:id", partController.GetPartByID)
}
