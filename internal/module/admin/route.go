package admin

import (
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(router fiber.Router, db *db.Queries, cfg *config.Config, jwtAuth fiber.Handler) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))
	router.Use(jwtAuth)
	router.Use(middleware.RequireRoles(auth.RoleSuperAdmin, auth.RoleStaff))

	repo := repository.NewAdminRepository(db)
	adminUsecase := usecase.NewAdminUsecase(repo)
	adminController := controller.NewAdminController(adminUsecase)

	router.Get("/", adminController.GetAllAdminUsers)
	router.Get("/:id", adminController.GetAdminByID)
	router.Put("/:id", adminController.UpdateAdmin)
	router.Delete("/:id", adminController.DeleteAdmin)
	router.Post("/register", adminController.CreateAdmin)
}

func SetupPublicRoutes(router fiber.Router, db *db.Queries, cfg *config.Config) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewAdminRepository(db)
	authUsecase := usecase.NewAuthUsecase(repo, cfg)
	authController := controller.NewAuthController(authUsecase, cfg)

	router.Post("/login", authController.Login)
}
