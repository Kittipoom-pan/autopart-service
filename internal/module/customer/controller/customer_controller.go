package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
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

func (h *CustomerController) GetCustomer(c *fiber.Ctx) error {
	idStr := c.Query("id")
	h.logger.Debug().Str("customer_id", idStr).Msg("Get customer request")
	if idStr == "" {
		h.logger.Warn().Msg("id parameter is missing")
		apiErr := customerror.InvalidRequestData(map[string]string{"id": "id is required"})
		return c.Status(apiErr.Code).JSON(common.BaseResponse{
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Result:  nil,
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("invalid id format")
		apiErr := customerror.InvalidRequestData(map[string]string{"id": "invalid id format"})
		return c.Status(apiErr.Code).JSON(common.BaseResponse{
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Result:  nil,
		})
	}

	customer, err := h.usecase.GetCustomerByID(c.Context(), id)
	if err != nil {
		if apiErr, ok := err.(customerror.APIError); ok {
			return c.Status(apiErr.Code).JSON(common.BaseResponse{
				Code:    apiErr.Code,
				Message: apiErr.Message,
				Result:  nil,
			})
		}

		if notFoundErr, ok := err.(*customerror.NotFoundError); ok {
			return c.Status(common.StatusNotFound).JSON(common.BaseResponse{
				Code:    common.StatusNotFound,
				Message: notFoundErr.Error(),
				Result:  nil,
			})
		}
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
	h.logger.Debug().Interface("customer_data", customer).Msg("Create customer request")
	if err := c.BodyParser(customer); err != nil {
		apiErr := customerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"})
		return c.Status(apiErr.Code).JSON(common.BaseResponse{
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Result:  nil,
		})
	}

	customerID, err := h.usecase.CreateCustomer(c.Context(), customer)
	if err != nil {
		apiErr := customerror.NewAPIError(common.StatusError, err.Error())
		return c.Status(apiErr.Code).JSON(common.BaseResponse{
			Code:    apiErr.Code,
			Message: apiErr.Message,
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
