package order

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func SetupPublicRoutes(
	router fiber.Router,
	conn *sql.DB,
	cfg *config.Config,
	logger zerolog.Logger,
) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	repo := repository.NewOrderRepository(conn, logger)
	orderUsecase := usecase.NewOrderUsecase(repo, logger)
	orderController := controller.NewOrderController(orderUsecase, logger)

	router.Get("/", orderController.GetAllOrders)
	router.Get("/:id", orderController.GetOrderByID)
	router.Get("/order-number/:orderNumber", orderController.GetOrderByOrderNumber)
}

func SetupPrivateRoutes(
	router fiber.Router,
	conn *sql.DB,
	cfg *config.Config,
	jwtAuth fiber.Handler,
	logger zerolog.Logger,
) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))
	router.Use(jwtAuth)

	repo := repository.NewOrderRepository(conn, logger)
	orderUsecase := usecase.NewOrderUsecase(repo, logger)
	orderController := controller.NewOrderController(orderUsecase, logger)

	router.Post("/", orderController.CreateOrderByCustomer)
	router.Post("/admin", orderController.CreateOrderByAdmin)
	router.Get("/", orderController.GetAllOrders)
	router.Get("/:id", orderController.GetOrderByID)
	router.Get("/order-number/:orderNumber", orderController.GetOrderByOrderNumber)
}
