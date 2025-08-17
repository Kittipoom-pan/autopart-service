package admin

import (
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/usecase"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router, db *db.Queries) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	// create dependencies
	repo := repository.NewAdminRepository(db)
	usecase := usecase.NewAdminUsecase(repo)
	controller := controller.NewAdminController(usecase)

	router.Get("/", controller.GetAllAdminUsers)
	router.Get("/:id", controller.GetAdminByID)
	router.Post("/", controller.CreateAdmin)
	router.Put("/:id", controller.UpdateAdmin)
	router.Delete("/:id", controller.DeleteAdmin)
}
