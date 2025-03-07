package controller

import (
	"fmt"
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase"
	"github.com/gofiber/fiber/v2"
)

type CustomerController interface {
	GetCustomer(c *fiber.Ctx) error
	CreateCustomer(c *fiber.Ctx) error
}

type CustomerControllerImpl struct {
	usecase usecase.CustomerUsecase
}

func NewCustomerController(usecase usecase.CustomerUsecase) CustomerController {
	return &CustomerControllerImpl{usecase: usecase}
}

func (h *CustomerControllerImpl) GetCustomer(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format"})
	}

	user, err := h.usecase.GetCustomerByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if user == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(user)
}

func (h *CustomerControllerImpl) CreateCustomer(c *fiber.Ctx) error {
	user := new(entitie.Customer)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println("User data:", user)
	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	customerID, err := h.usecase.CreateCustomer(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"customer_id": customerID})
}
