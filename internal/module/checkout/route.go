package checkout

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	addressRepository "github.com/Kittipoom-pan/autopart-service/internal/module/address/repository"
	cartRepository "github.com/Kittipoom-pan/autopart-service/internal/module/cart/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout/usecase"
	partRepository "github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	pricingUsecase "github.com/Kittipoom-pan/autopart-service/internal/module/pricing/usecase"
	shippingRepository "github.com/Kittipoom-pan/autopart-service/internal/module/shipping/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

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

	cartRepo := cartRepository.NewCartRepository(conn, logger)
	partRepo := partRepository.NewPartRepository(db, logger)
	addressRepo := addressRepository.NewAddressRepository(conn, logger)
	shippingRateRepo := shippingRepository.NewShippingRateRepository(conn, logger)
	pricingUsecase := pricingUsecase.NewPricingUsecase(shippingRateRepo, logger)

	checkoutUsecase := usecase.NewCheckoutUsecase(cartRepo, partRepo, addressRepo, pricingUsecase, logger)
	checkoutController := controller.NewCheckoutController(checkoutUsecase, logger)

	router.Post("/preview", checkoutController.CheckoutPreview)
}
