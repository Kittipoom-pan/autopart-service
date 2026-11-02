package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	_ "github.com/Kittipoom-pan/autopart-service/internal/common"
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

// GetCart godoc
// @Summary      Get cart
// @Description  ดึงตะกร้าของลูกค้า
// @Tags         cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  common.BaseResponse[entity.CartRes]
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/cart [get]
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

// AddItem godoc
// @Summary      Add item to cart
// @Description  เพิ่มสินค้าเข้าตะกร้า
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      entity.AddCartItemReq  true  "Cart item data"
// @Success      200   {object}  common.BaseResponse[entity.CartRes]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/cart/items [post]
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

// UpdateItemQuantity godoc
// @Summary      Update cart item quantity
// @Description  แก้ไขจำนวนสินค้าในตะกร้า
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                               true  "Cart Item ID"
// @Param        body  body      entity.UpdateCartItemQuantityReq  true  "Cart item quantity data"
// @Success      200   {object}  common.BaseResponse[entity.CartRes]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      404   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/cart/items/{id} [patch]
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
