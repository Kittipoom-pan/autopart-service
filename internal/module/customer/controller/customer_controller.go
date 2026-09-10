package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entity"
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

	customers, err := h.usecase.GetAllCustomers(c.UserContext())
	if err != nil {
		return helper.RespondError(c, err)
	}
	if customers == nil {
		customers = []*entity.CustomerRes{}
	}

	return helper.RespondSuccess(c, fiber.StatusOK, customers, "Customers retrieved successfully")
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

	if err := auth.EnsureSelfOrAdmin(c, id, h.logger); err != nil {
		return helper.RespondError(c, err)
	}

	customer, err := h.usecase.GetCustomerByID(c.UserContext(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, customer, "")
}

func (h *CustomerController) CreateCustomer(c *fiber.Ctx) error {
	customer := new(entity.CustomerReq)
	if err := c.BodyParser(customer); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	customerID, err := h.usecase.CreateCustomer(c.UserContext(), customer)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, fiber.Map{"customer_id": customerID}, "Customer created successfully")
}

func (h *CustomerController) UpdateCustomer(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	idStr := c.Params("id")
	customerID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid customer ID format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid customer ID format"}))
	}

	if err := auth.EnsureSelfOrAdmin(c, customerID, h.logger); err != nil {
		return helper.RespondError(c, err)
	}

	customerReq := new(entity.CustomerReq)
	if err := c.BodyParser(customerReq); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse request body")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	if err := h.usecase.UpdateCustomer(c.UserContext(), customerID, customerReq, userID); err != nil {
		h.logger.Error().Err(err).Int("customer_id", customerID).Msg("Failed to update customer")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("customer_id", customerID).Msg("Customer updated successfully")
	return helper.RespondSuccess(c, fiber.StatusOK, nil, "Customer updated successfully")
}

func (h *CustomerController) DeleteCustomer(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	idStr := c.Params("id")
	customerID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Str("id_param", idStr).Msg("Invalid customer ID format")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"id": "invalid customer ID format"}))
	}

	if err := auth.EnsureSelfOrAdmin(c, customerID, h.logger); err != nil {
		return helper.RespondError(c, err)
	}

	if err := h.usecase.DeleteCustomer(c.UserContext(), customerID, userID); err != nil {
		h.logger.Error().Err(err).Int("customer_id", customerID).Msg("Failed to delete customer")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("customer_id", customerID).Int("by_user", userID).Msg("Customer deleted successfully")
	return helper.RespondSuccess(c, fiber.StatusOK, nil, "Customer deleted successfully")
}
