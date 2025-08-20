package admin

import (
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, auth fiber.Handler) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	// create dependencies
	repo := repository.NewAdminRepository(db)
	usecase := usecase.NewAdminUsecase(repo)
	controller := controller.NewAdminController(usecase)

	router.Get("/", auth, controller.GetAllAdminUsers)
	router.Get("/:id", auth, controller.GetAdminByID)
	router.Put("/:id", auth, controller.UpdateAdmin)
	router.Delete("/:id", auth, controller.DeleteAdmin)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewAdminRepository(db)
	authUsecase := usecase.NewAuthUsecase(repo, cfg)
	usecase := usecase.NewAdminUsecase(repo)
	adminController := controller.NewAdminController(usecase)
	authController := controller.NewAuthController(authUsecase, cfg)

	router.Post("/login", authController.Login)
	router.Post("/register", adminController.CreateAdmin)
}
