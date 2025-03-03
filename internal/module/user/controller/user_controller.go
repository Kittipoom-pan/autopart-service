package controller

import (
	"github.com/Kittipoom-pan/autopart-service/internal/module/user/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/user/usecase"
	"github.com/gofiber/fiber/v2"
)

type UsersController interface {
	GetUsers(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	// ... other methods ...
}

type UsersControllerImpl struct {
	usecase usecase.UsersUsecase
}

func NewUsersController(usecase usecase.UsersUsecase) UsersController {
	return &UsersControllerImpl{usecase: usecase}
}

func (h *UsersControllerImpl) GetUsers(c *fiber.Ctx) error {
	// ... implementation ...
	return nil
}

func (h *UsersControllerImpl) CreateUser(c *fiber.Ctx) error {
	user := new(entitie.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}
	if err := h.usecase.CreateUser(user); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusCreated)
}
