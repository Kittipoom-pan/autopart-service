package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CustomerController struct {
	usecase usecase.CustomerUsecase
	logger  zerolog.Logger
}

func NewCustomerController(usecase usecase.CustomerUsecase) *CustomerController {
	return &CustomerController{
		usecase: usecase,
		logger:  log.With().Str("component", "customer_controller").Logger(),
	}
}

func (h *CustomerController) GetAllCustomers(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get all customers request")

	customers, err := h.usecase.GetAllCustomers(c.Context())
	if err != nil {
		return helper.RespondError(c, err)
	}
	if len(customers) == 0 {
		return helper.RespondSuccess(c, fiber.StatusOK, customers, "Customers not found")
	}

	return helper.RespondSuccess(c, fiber.StatusOK, customers, "")
}

func (h *CustomerController) GetCustomerByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("customer_id", idStr).Msg("Get customer request")

	if idStr == "" {
		h.logger.Warn().Msg("Id parameter is missing")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "id is required"}))
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn().Msg("Invalid id format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid id format"}))
	}

	customer, err := h.usecase.GetCustomerByID(c.Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, customer, "")
}

func (h *CustomerController) CreateCustomer(c *fiber.Ctx) error {
	customer := new(entitie.CustomerReq)
	if err := c.BodyParser(customer); err != nil {
		h.logger.Warn().Err(err).Bytes("raw_body", c.Body()).Msg("Failed to parse request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	customerID, err := h.usecase.CreateCustomer(c.Context(), customer)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, fiber.Map{"customer_id": customerID}, "Customer created successfully")
}

func (h *CustomerController) UpdateCustomer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid customer ID format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid customer ID format"}))
	}

	customerReq := new(entitie.CustomerReq)
	if err := c.BodyParser(customerReq); err != nil {
		h.logger.Warn().Err(err).Bytes("raw_body", c.Body()).Msg("Failed to parse request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	if err := h.usecase.UpdateCustomer(c.Context(), id, customerReq); err != nil {
		h.logger.Error().Err(err).Int("customer_id", id).Msg("Failed to update customer")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("customer_id", id).Msg("Customer updated successfully")
	return helper.RespondSuccess(c, fiber.StatusOK, nil, "Customer updated successfully")
}

func (h *CustomerController) DeleteCustomer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid customer ID format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid customer ID format"}))
	}

	if err := h.usecase.DeleteCustomer(c.Context(), id); err != nil {
		h.logger.Error().Err(err).Int("customer_id", id).Msg("Failed to delete customer")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("customer_id", id).Msg("Customer deleted successfully")
	return helper.RespondSuccess(c, fiber.StatusOK, nil, "Customer deleted successfully")
}
