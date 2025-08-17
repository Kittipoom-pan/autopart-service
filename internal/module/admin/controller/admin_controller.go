package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/usecase"
	adminerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AdminController struct {
	usecase usecase.AdminUsecase
	logger  zerolog.Logger
}

func NewAdminController(usecase usecase.AdminUsecase) *AdminController {
	return &AdminController{
		usecase: usecase,
		logger:  log.With().Str("component", "admin_controller").Logger(),
	}
}

func (h *AdminController) GetAllAdminUsers(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get all admins request")

	admins, err := h.usecase.GetAllAdmins(c.Context())

	if err != nil {
		return helper.RespondError(c, err)
	}
	if len(admins) == 0 {
		return helper.RespondSuccess(c, fiber.StatusOK, admins, "Admins not found")
	}

	return helper.RespondSuccess(c, fiber.StatusOK, admins, "")
}

func (h *AdminController) GetAdminByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("admin_id", idStr).Msg("Get admin request")

	if idStr == "" {
		h.logger.Warn().Msg("Id parameter is missing")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"id": "id is required"}))
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn().Msg("Invalid id format")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"id": "invalid id format"}))
	}

	admin, err := h.usecase.GetAdminByID(c.Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, admin, "")
}

func (h *AdminController) CreateAdmin(c *fiber.Ctx) error {
	admin := new(entitie.AdminReq)
	if err := c.BodyParser(admin); err != nil {
		h.logger.Warn().Err(err).Bytes("raw_body", c.Body()).Msg("Failed to parse request body")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	adminID, err := h.usecase.CreateAdmin(c.Context(), admin)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, fiber.Map{"admin_id": adminID}, "Admin created successfully")
}

func (h *AdminController) UpdateAdmin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid admin ID format")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"id": "invalid admin ID format"}))
	}

	adminReq := new(entitie.AdminReq)
	if err := c.BodyParser(adminReq); err != nil {
		h.logger.Warn().Err(err).Bytes("raw_body", c.Body()).Msg("Failed to parse request body")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	if err := h.usecase.UpdateAdmin(c.Context(), id, adminReq); err != nil {
		h.logger.Error().Err(err).Int("admin_id", id).Msg("Failed to update admin")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("admin_id", id).Msg("Admin updated successfully")
	return helper.RespondSuccess(c, fiber.StatusOK, nil, "Admin updated successfully")
}

func (h *AdminController) DeleteAdmin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid admin ID format")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"id": "invalid admin ID format"}))
	}

	if err := h.usecase.DeleteAdmin(c.Context(), id); err != nil {
		h.logger.Error().Err(err).Int("admin_id", id).Msg("Failed to delete admin")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("admin_id", id).Msg("Admin deleted successfully")
	return helper.RespondSuccess(c, fiber.StatusOK, nil, "Admin deleted successfully")
}
