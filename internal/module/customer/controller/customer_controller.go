package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase"
	"github.com/gofiber/fiber/v2"
)

type CustomerController struct {
	usecase usecase.CustomerUsecase
}

func NewCustomerController(usecase usecase.CustomerUsecase) *CustomerController {
	return &CustomerController{usecase: usecase}
}

func (h *CustomerController) GetCustomer(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse{
			Code:    common.StatusBadRequest,
			Message: "id is required",
			Result:  nil,
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.BaseResponse{
			Code:    common.StatusBadRequest,
			Message: "invalid id format",
			Result:  nil,
		})
	}

	customer, err := h.usecase.GetCustomerByID(c.Context(), id)
	if err != nil {
		return c.Status(common.StatusError).JSON(common.BaseResponse{
			Code:    common.StatusError,
			Message: err.Error(),
			Result:  nil,
		})
	}

	if customer == nil {
		return c.Status(common.StatusNotFound).JSON(common.BaseResponse{
			Code:    common.StatusNotFound,
			Message: "Customer not found",
			Result:  nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(common.BaseResponse{
		Code:    common.StatusSuccess,
		Message: "success",
		Result:  customer,
	})
}

func (h *CustomerController) CreateCustomer(c *fiber.Ctx) error {
	customer := new(entitie.CustomerReq)
	if err := c.BodyParser(customer); err != nil {
		return c.Status(common.StatusError).JSON(common.BaseResponse{
			Code:    common.StatusError,
			Message: err.Error(),
			Result:  nil,
		})
	}

	customerID, err := h.usecase.CreateCustomer(c.Context(), customer)
	if err != nil {
		return c.Status(common.StatusError).JSON(common.BaseResponse{
			Code:    common.StatusError,
			Message: err.Error(),
			Result:  nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(common.BaseResponse{
		Code:    common.StatusSuccess,
		Message: "success",
		Result:  fiber.Map{"customer_id": customerID},
	})
}
