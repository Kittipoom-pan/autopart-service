package customer_payment_method

import (
	"database/sql"
	"time"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/middleware"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/controller"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/usecase"
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
	repo := repository.NewCustomerPaymentMethodRepository(conn, logger)
	paymentUsecase := usecase.NewCustomerPaymentMethodUsecase(repo, txManager, logger)
	paymentController := controller.NewCustomerPaymentMethodController(paymentUsecase, logger)

	router.Get("/", paymentController.GetPaymentMethods)
	router.Get("/:id", paymentController.GetPaymentMethodByID)
	router.Post("/", paymentController.CreatePaymentMethod)
	router.Put("/:id", paymentController.UpdatePaymentMethod)
	router.Delete("/:id", paymentController.DeletePaymentMethod)
	router.Patch("/:id/default", paymentController.SetDefaultPaymentMethod)
}
