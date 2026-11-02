package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	_ "github.com/Kittipoom-pan/autopart-service/internal/common"
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

// GetPaymentMethods godoc
// @Summary      List payment methods
// @Description  ดึงรายการบัตร/ช่องทางจ่ายที่บันทึกไว้
// @Tags         payment-method
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  common.BaseResponse[[]entity.CustomerPaymentMethodRes]
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/payment-method [get]
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

// GetPaymentMethodByID godoc
// @Summary      Get payment method
// @Description  ดึงช่องทางจ่ายตาม ID
// @Tags         payment-method
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Payment Method ID"
// @Success      200  {object}  common.BaseResponse[entity.CustomerPaymentMethodRes]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/payment-method/{id} [get]
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

// CreatePaymentMethod godoc
// @Summary      Create payment method
// @Description  เพิ่มช่องทางจ่าย
// @Tags         payment-method
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      entity.CreateCustomerPaymentMethodReq  true  "Payment method data"
// @Success      201   {object}  common.BaseResponse[entity.CustomerPaymentMethodRes]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/payment-method [post]
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

// UpdatePaymentMethod godoc
// @Summary      Update payment method
// @Description  แก้ไขช่องทางจ่าย
// @Tags         payment-method
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                                   true  "Payment Method ID"
// @Param        body  body      entity.UpdateCustomerPaymentMethodReq  true  "Payment method data"
// @Success      200   {object}  common.BaseResponse[entity.CustomerPaymentMethodRes]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      404   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/payment-method/{id} [put]
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

// DeletePaymentMethod godoc
// @Summary      Delete payment method
// @Description  ลบช่องทางจ่าย
// @Tags         payment-method
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Payment Method ID"
// @Success      200  {object}  common.BaseResponse[any]
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/payment-method/{id} [delete]
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

// SetDefaultPaymentMethod godoc
// @Summary      Set default payment method
// @Description  ตั้งช่องทางจ่ายเริ่มต้น
// @Tags         payment-method
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Payment Method ID"
// @Success      200  {object}  common.BaseResponse[entity.CustomerPaymentMethodRes]
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/payment-method/{id}/default [patch]
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
