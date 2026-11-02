package address

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/usecase"
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

	txManager := database.NewTxManager(conn, logger)
	repo := repository.NewAddressRepository(conn, logger)
	addressUsecase := usecase.NewAddressUsecase(repo, txManager, logger)
	addressController := controller.NewAddressController(addressUsecase, logger)

	router.Get("/", addressController.GetAddresses)
	router.Get("/:id", addressController.GetAddressByIDAndCustomer)
	router.Post("/", addressController.CreateAddress)
	router.Put("/:id", addressController.UpdateAddress)
	router.Delete("/:id", addressController.DeleteAddress)
	router.Patch("/:id/default", addressController.SetDefaultAddress)
}
