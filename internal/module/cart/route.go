package cart

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/usecase"
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
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))
	router.Use(jwtAuth)

	repo := repository.NewCartRepository(conn, logger)
	cartUsecase := usecase.NewCartUsecase(repo, logger)
	cartController := controller.NewCartController(cartUsecase, logger)

	router.Get("/", cartController.GetCart)
	router.Post("/items", cartController.AddItem)
	router.Patch("/items/:id", cartController.UpdateItemQuantity)
}
