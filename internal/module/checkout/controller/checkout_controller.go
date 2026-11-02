package controller

import (
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout/usecase"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type CheckoutController struct {
	usecase usecase.CheckoutUsecase
	logger  zerolog.Logger
}

func NewCheckoutController(usecase usecase.CheckoutUsecase, logger zerolog.Logger) *CheckoutController {
	return &CheckoutController{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *CheckoutController) CheckoutPreview(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get checkout preview request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	var req entity.CheckoutPreviewReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	checkout, err := h.usecase.PreviewCheckout(c.UserContext(), int32(customerID), &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, checkout, "Checkout preview retrieved successfully")
}
