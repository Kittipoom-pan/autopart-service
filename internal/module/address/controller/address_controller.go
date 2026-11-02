package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
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
