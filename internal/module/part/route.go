package part

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	imagerepository "github.com/Kittipoom-pan/autopart-service/internal/module/image/repository"
	imageusecase "github.com/Kittipoom-pan/autopart-service/internal/module/image/usecase"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/usecase"
	"github.com/Kittipoom-pan/autopart-service/pkg/imagestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func SetupPrivateRoutes(router fiber.Router, conn *sql.DB, cfg *config.Config, jwtAuth fiber.Handler, logger zerolog.Logger) {
	router.Use(middleware.TimeoutMiddleware(30 * time.Second))
	router.Use(jwtAuth)
	router.Use(middleware.RequireRoles(auth.RoleSuperAdmin, auth.RoleStaff))

	store := imagestore.NewLocalStore(cfg.Upload.Dir, cfg.Upload.BaseURL, cfg.Upload.MaxWidth, logger)
	imageRepo := imagerepository.NewImageRepository(conn, logger)
	txManager := database.NewTxManager(conn, logger)
	imageUsecase := imageusecase.NewImageUsecase(
		imageRepo, store, txManager,
		cfg.Upload.MaxSizeBytes, cfg.Upload.MaxFiles, logger,
	)

	repo := repository.NewPartRepository(db.New(conn), logger)
	partUsecase := usecase.NewPartUsecase(repo, imageUsecase, txManager, logger)
	partController := controller.NewPartController(partUsecase, logger)

	router.Post("/", partController.CreatePart)
	router.Put("/:id", partController.UpdatePart)
	router.Delete("/:id", partController.DeletePart)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, logger zerolog.Logger) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewPartRepository(db, logger)
	partUsecase := usecase.NewPartUsecase(repo, nil, nil, logger)
	partController := controller.NewPartController(partUsecase, logger)

	router.Get("/", partController.GetAllParts)
	router.Get("/:id", partController.GetPartByID)
}
