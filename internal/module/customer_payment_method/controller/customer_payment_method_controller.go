package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/usecase"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type CustomerPaymentMethodController struct {
	usecase usecase.CustomerPaymentMethodUsecase
	logger  zerolog.Logger
}

func NewCustomerPaymentMethodController(usecase usecase.CustomerPaymentMethodUsecase, logger zerolog.Logger) *CustomerPaymentMethodController {
	return &CustomerPaymentMethodController{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *CustomerPaymentMethodController) GetPaymentMethods(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get customer payment methods request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	methods, err := h.usecase.GetCustomerPaymentMethods(c.UserContext(), int32(customerID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, methods, "Payment methods retrieved successfully")
}

func (h *CustomerPaymentMethodController) GetPaymentMethodByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("payment_method_id", idStr).Msg("Get customer payment method by id request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid payment method id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid payment method id format"}))
	}

	method, err := h.usecase.GetCustomerPaymentMethodByID(c.UserContext(), int32(id), int32(customerID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, method, "Payment method retrieved successfully")
}

func (h *CustomerPaymentMethodController) CreatePaymentMethod(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Create customer payment method request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	var req entity.CreateCustomerPaymentMethodReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	method, err := h.usecase.CreateCustomerPaymentMethod(c.UserContext(), int32(customerID), &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, method, "Payment method created successfully")
}

func (h *CustomerPaymentMethodController) UpdatePaymentMethod(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("payment_method_id", idStr).Msg("Update customer payment method request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid payment method id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid payment method id format"}))
	}

	var req entity.UpdateCustomerPaymentMethodReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	method, err := h.usecase.UpdateCustomerPaymentMethod(c.UserContext(), int32(id), int32(customerID), &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, method, "Payment method updated successfully")
}

func (h *CustomerPaymentMethodController) DeletePaymentMethod(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("payment_method_id", idStr).Msg("Delete customer payment method request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid payment method id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid payment method id format"}))
	}

	if err := h.usecase.DeleteCustomerPaymentMethod(c.UserContext(), int32(id), int32(customerID)); err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, map[string]string{"message": "Payment method deleted successfully"}, "Payment method deleted successfully")
}

func (h *CustomerPaymentMethodController) SetDefaultPaymentMethod(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("payment_method_id", idStr).Msg("Set default customer payment method request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid payment method id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid payment method id format"}))
	}

	method, err := h.usecase.SetDefaultCustomerPaymentMethod(c.UserContext(), int32(id), int32(customerID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, method, "Default payment method set successfully")
}
