package part

import (
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func SetupPrivateRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, jwtAuth fiber.Handler, logger zerolog.Logger) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))
	router.Use(jwtAuth)
	router.Use(middleware.RequireRoles(auth.RoleSuperAdmin, auth.RoleStaff, auth.RoleCustomer))

	repo := repository.NewPartRepository(db, logger)
	partUsecase := usecase.NewPartUsecase(repo, logger)
	partController := controller.NewPartController(partUsecase, logger)

	router.Post("/", partController.CreatePart)
	router.Put("/:id", partController.UpdatePart)
	router.Delete("/:id", partController.DeletePart)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, logger zerolog.Logger) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewPartRepository(db, logger)
	partUsecase := usecase.NewPartUsecase(repo, logger)
	partController := controller.NewPartController(partUsecase, logger)

	router.Get("/", partController.GetAllParts)
	router.Get("/:id", partController.GetPartByID)
}
