package image

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/usecase"
	"github.com/Kittipoom-pan/autopart-service/pkg/imagestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func SetupPrivateRoutes(
	router fiber.Router,
	conn *sql.DB,
	cfg *config.Config,
	jwtAuth fiber.Handler,
	logger zerolog.Logger,
) {
	router.Use(middleware.TimeoutMiddleware(30 * time.Second))
	router.Use(jwtAuth)
	router.Use(middleware.RequireRoles(auth.RoleSuperAdmin, auth.RoleStaff))

	store := imagestore.NewLocalStore(cfg.Upload.Dir, cfg.Upload.BaseURL, cfg.Upload.MaxWidth, logger)
	repo := repository.NewImageRepository(conn, logger)
	txManager := database.NewTxManager(conn, logger)
	imageUsecase := usecase.NewImageUsecase(
		repo, store, txManager,
		cfg.Upload.MaxSizeBytes, cfg.Upload.MaxFiles, logger,
	)
	imageController := controller.NewImageController(imageUsecase, logger)

	router.Post("/", imageController.UploadImage)
	router.Post("/multiple", imageController.UploadImages)
}
