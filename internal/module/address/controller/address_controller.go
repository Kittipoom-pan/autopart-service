package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	_ "github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/usecase"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AddressController struct {
	usecase usecase.AddressUsecase
	logger  zerolog.Logger
}

func NewAddressController(usecase usecase.AddressUsecase, logger zerolog.Logger) *AddressController {
	return &AddressController{
		usecase: usecase,
		logger:  logger,
	}
}

// GetAddresses godoc
// @Summary      List addresses
// @Description  ดึงรายการที่อยู่ของลูกค้า
// @Tags         address
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  common.BaseResponse[[]entity.AddressRes]
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/address [get]
func (h *AddressController) GetAddresses(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get customer addresses request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	addresses, err := h.usecase.GetCustomerAddresses(c.UserContext(), int32(customerID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, addresses, "Addresses retrieved successfully")
}

// GetAddressByIDAndCustomer godoc
// @Summary      Get address
// @Description  ดึงที่อยู่ตาม ID
// @Tags         address
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Address ID"
// @Success      200  {object}  common.BaseResponse[entity.AddressRes]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/address/{id} [get]
func (h *AddressController) GetAddressByIDAndCustomer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("address_id", idStr).Msg("Get address by id request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid address id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid address id format"}))
	}

	address, err := h.usecase.GetAddressByIDAndCustomer(c.UserContext(), int32(id), int32(customerID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, address, "Address retrieved successfully")
}

// CreateAddress godoc
// @Summary      Create address
// @Description  เพิ่มที่อยู่ใหม่
// @Tags         address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      entity.CreateAddressReq  true  "Address data"
// @Success      201   {object}  common.BaseResponse[entity.AddressRes]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/address [post]
func (h *AddressController) CreateAddress(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Create address request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	var req entity.CreateAddressReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	address, err := h.usecase.CreateAddress(c.UserContext(), int32(customerID), &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, address, "Address created successfully")
}

// UpdateAddress godoc
// @Summary      Update address
// @Description  แก้ไขที่อยู่
// @Tags         address
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                      true  "Address ID"
// @Param        body  body      entity.UpdateAddressReq  true  "Address data"
// @Success      200   {object}  common.BaseResponse[entity.AddressRes]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      404   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/address/{id} [put]
func (h *AddressController) UpdateAddress(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("address_id", idStr).Msg("Update address request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid address id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid address id format"}))
	}

	var req entity.UpdateAddressReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid request body"}))
	}

	address, err := h.usecase.UpdateAddress(c.UserContext(), int32(id), int32(customerID), &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, address, "Address updated successfully")
}

// DeleteAddress godoc
// @Summary      Delete address
// @Description  ลบที่อยู่
// @Tags         address
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Address ID"
// @Success      200  {object}  common.BaseResponse[any]
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/address/{id} [delete]
func (h *AddressController) DeleteAddress(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("address_id", idStr).Msg("Delete address request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid address id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid address id format"}))
	}

	if err := h.usecase.DeleteAddress(c.UserContext(), int32(id), int32(customerID)); err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, map[string]string{"message": "Address deleted successfully"}, "Address deleted successfully")
}

// SetDefaultAddress godoc
// @Summary      Set default address
// @Description  ตั้งที่อยู่เริ่มต้น
// @Tags         address
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Address ID"
// @Success      200  {object}  common.BaseResponse[entity.AddressRes]
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/address/{id}/default [patch]
func (h *AddressController) SetDefaultAddress(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("address_id", idStr).Msg("Set default address request")

	customerID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.logger.Warn().Err(err).Msg("Invalid address id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid address id format"}))
	}

	address, err := h.usecase.SetDefaultAddress(c.UserContext(), int32(id), int32(customerID))
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, address, "Default address set successfully")
}
