package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/usecase"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type CartController struct {
	usecase usecase.CartUsecase
	logger  zerolog.Logger
}

func NewCartController(usecase usecase.CartUsecase, logger zerolog.Logger) *CartController {
	return &CartController{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *CartController) GetCart(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get cart request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	cart, err := h.usecase.GetCart(c.UserContext(), int32(customerID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, cart, "Cart retrieved successfully")
}

func (h *CartController) AddItem(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Add item to cart request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	var req entity.AddCartItemReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	cart, err := h.usecase.AddItem(c.UserContext(), int32(customerID), &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, cart, "Item added to cart successfully")
}

func (h *CartController) UpdateItemQuantity(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("cart_item_id", idStr).Msg("Update cart item quantity request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	cartItemID, err := strconv.Atoi(idStr)
	if err != nil || cartItemID <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid cart item id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid cart item id format"}))
	}

	var req entity.UpdateCartItemQuantityReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	cart, err := h.usecase.UpdateItemQuantity(c.UserContext(), int32(customerID), int32(cartItemID), &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, cart, "Cart item updated successfully")
}
