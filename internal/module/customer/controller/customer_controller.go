package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase"
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
		logger:  log.With().Str("component", "customer_controller").Logger(), // สร้าง logger instance
	}
}

func (h *CustomerController) GetCustomer(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		h.logger.Warn().Msg("id parameter is missing")
		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse{
			Code:    common.StatusBadRequest,
			Message: "id is required",
			Result:  nil,
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("invalid id format")
		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse{
			Code:    common.StatusBadRequest,
			Message: "invalid id format",
			Result:  nil,
		})
	}

	customer, err := h.usecase.GetCustomerByID(c.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get customer")
		return c.Status(common.StatusError).JSON(common.BaseResponse{
			Code:    common.StatusError,
			Message: "Failed to get customer: " + err.Error(),
			Result:  nil,
		})
	}

	if customer == nil {
		h.logger.Info().Int("customer_id", id).Msg("customer not found")
		return c.Status(common.StatusNotFound).JSON(common.BaseResponse{
			Code:    common.StatusNotFound,
			Message: "Customer not found",
			Result:  nil,
		})
	}

	h.logger.Info().Int("customer_id", id).Msg("customer retrieved successfully")
	return c.Status(fiber.StatusOK).JSON(common.BaseResponse{
		Code:    common.StatusSuccess,
		Message: "",
		Result:  customer,
	})
}

func (h *CustomerController) CreateCustomer(c *fiber.Ctx) error {
	customer := new(entitie.CustomerReq)
	if err := c.BodyParser(customer); err != nil {
		h.logger.Error().Err(err).Msg("failed to parse request body")
		return c.Status(common.StatusError).JSON(common.BaseResponse{
			Code:    common.StatusError,
			Message: "Invalid request body: " + err.Error(),
			Result:  nil,
		})
	}

	customerID, err := h.usecase.CreateCustomer(c.Context(), customer)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create customer")
		return c.Status(common.StatusError).JSON(common.BaseResponse{
			Code:    common.StatusError,
			Message: err.Error(),
			Result:  nil,
		})

	}
	h.logger.Info().Int64("customer_id", customerID).Msg("customer created successfully")
	return c.Status(fiber.StatusCreated).JSON(common.BaseResponse{
		Code:    common.StatusSuccess,
		Message: "",
		Result:  fiber.Map{"customer_id": customerID},
	})
}
