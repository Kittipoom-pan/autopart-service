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
)

func SetupPrivateRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, jwtAuth fiber.Handler) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))
	router.Use(jwtAuth)
	router.Use(middleware.RequireRoles(auth.RoleSuperAdmin, auth.RoleStaff, auth.RoleCustomer))

	repo := repository.NewPartRepository(db)
	partUsecase := usecase.NewPartUsecase(repo)
	partController := controller.NewPartController(partUsecase)

	router.Post("/", partController.CreatePart)
	router.Put("/:id", partController.UpdatePart)
	router.Delete("/:id", partController.DeletePart)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewPartRepository(db)
	partUsecase := usecase.NewPartUsecase(repo)
	partController := controller.NewPartController(partUsecase)

	router.Get("/", partController.GetAllParts)
	router.Get("/:id", partController.GetPartByID)
}
