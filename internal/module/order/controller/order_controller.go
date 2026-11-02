package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/usecase"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type OrderController struct {
	usecase usecase.OrderUsecase
	logger  zerolog.Logger
}

func NewOrderController(usecase usecase.OrderUsecase, logger zerolog.Logger) *OrderController {
	return &OrderController{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *OrderController) GetAllOrders(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get all orders request")

	orders, err := h.usecase.GetAllOrders(c.UserContext())
	if err != nil {
		return helper.RespondError(c, err)
	}

	if orders == nil {
		orders = []*entity.OrderRes{}
	}

	return helper.RespondSuccess(c, fiber.StatusOK, orders, "Orders retrieved successfully")
}

func (h *OrderController) GetOrderByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("order_id", idStr).Msg("Get order request")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Invalid order id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid order id format"}))
	}

	order, err := h.usecase.GetOrderByID(c.UserContext(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, order, "")
}

func (h *OrderController) GetOrderByOrderNumber(c *fiber.Ctx) error {
	orderNumber := c.Params("orderNumber")
	h.logger.Debug().Str("order_number", orderNumber).Msg("Get order by order number request")

	if orderNumber == "" {
		h.logger.Warn().Msg("Order number parameter is missing")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"order_number": "order_number is required"}))
	}

	order, err := h.usecase.GetOrderByOrderNumber(c.UserContext(), orderNumber)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, order, "")
}

func (h *OrderController) CreateOrderByCustomer(c *fiber.Ctx) error {
	var req entity.CustomerCreateOrderReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	customerID, err := auth.GetUserID(c, h.logger) // ต้องเป็น middleware ของลูกค้า
	if err != nil {
		return helper.RespondError(c, err)
	}

	orderNo, err := h.usecase.CreateOrderByCustomer(c.UserContext(), &req, int32(customerID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, map[string]string{"order_number": orderNo}, "Order created successfully")
}

func (h *OrderController) CreateOrderByAdmin(c *fiber.Ctx) error {
	var req entity.AdminCreateOrderReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	adminID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	orderNo, err := h.usecase.CreateOrderByAdmin(c.UserContext(), &req, int32(adminID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, map[string]string{"order_number": orderNo}, "Order created successfully")
}
