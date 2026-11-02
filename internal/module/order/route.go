package order

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	addressRepository "github.com/Kittipoom-pan/autopart-service/internal/module/address/repository"
	cartRepository "github.com/Kittipoom-pan/autopart-service/internal/module/cart/repository"
	cpmRepository "github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/usecase"
	partRepository "github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	paymentRepository "github.com/Kittipoom-pan/autopart-service/internal/module/payment/repository"
	paymentUsecase "github.com/Kittipoom-pan/autopart-service/internal/module/payment/usecase"
	pricingUsecase "github.com/Kittipoom-pan/autopart-service/internal/module/pricing/usecase"
	shippingRepository "github.com/Kittipoom-pan/autopart-service/internal/module/shipping/repository"
	stockRepository "github.com/Kittipoom-pan/autopart-service/internal/module/stock/repository"
	stockUsecase "github.com/Kittipoom-pan/autopart-service/internal/module/stock/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func newOrderController(
	conn *sql.DB,
	db *db.Queries,
	logger zerolog.Logger,
) *controller.OrderController {
	txManager := database.NewTxManager(conn, logger)
	repo := repository.NewOrderRepository(conn, logger)
	cartRepo := cartRepository.NewCartRepository(conn, logger)
	partRepo := partRepository.NewPartRepository(db, logger)
	shippingRateRepo := shippingRepository.NewShippingRateRepository(conn, logger)
	pricingSvc := pricingUsecase.NewPricingUsecase(shippingRateRepo, logger)
	stockSvc := stockUsecase.NewStockUsecase(stockRepository.NewStockRepository(conn, logger), logger)
	paymentSvc := paymentUsecase.NewPaymentUsecase(
		paymentRepository.NewPaymentRepository(conn, logger),
		cpmRepository.NewCustomerPaymentMethodRepository(conn, logger),
		logger,
	)
	orderUsecase := usecase.NewOrderUsecase(
		repo,
		cartRepo,
		partRepo,
		addressRepository.NewAddressRepository(conn, logger),
		stockSvc,
		paymentSvc,
		txManager,
		pricingSvc,
		logger,
	)
	return controller.NewOrderController(orderUsecase, logger)
}

func SetupPublicRoutes(
	router fiber.Router,
	conn *sql.DB,
	db *db.Queries,
	cfg *config.Config,
	logger zerolog.Logger,
) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))

	orderController := newOrderController(conn, db, logger)

	router.Get("/", orderController.GetAllOrders)
	router.Get("/:id", orderController.GetOrderByID)
	router.Get("/order-number/:orderNumber", orderController.GetOrderByOrderNumber)
}

func SetupPrivateRoutes(
	router fiber.Router,
	conn *sql.DB,
	db *db.Queries,
	cfg *config.Config,
	jwtAuth fiber.Handler,
	logger zerolog.Logger,
) {
	router.Use(middleware.TimeoutMiddleware(3 * time.Second))
	router.Use(jwtAuth)

	orderController := newOrderController(conn, db, logger)

	router.Post("/", orderController.CreateOrderByCustomer)
	router.Post("/admin", orderController.CreateOrderByAdmin)
	router.Get("/", orderController.GetAllOrders)
	router.Get("/:id", orderController.GetOrderByID)
	router.Get("/order-number/:orderNumber", orderController.GetOrderByOrderNumber)
}
