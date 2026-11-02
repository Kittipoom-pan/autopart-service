package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	_ "github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/usecase"
	adminerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AdminController struct {
	usecase usecase.AdminUsecase
	logger  zerolog.Logger
}

func NewAdminController(usecase usecase.AdminUsecase, logger zerolog.Logger) *AdminController {
	return &AdminController{
		usecase: usecase,
		logger:  logger,
	}
}

// GetAllAdminUsers godoc
// @Summary      List admin users
// @Description  ดึงรายการผู้ดูแลระบบทั้งหมด
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  common.BaseResponse[[]entity.AdminRes]
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      403  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/admin [get]
func (h *AdminController) GetAllAdminUsers(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get all admins request")

	admins, err := h.usecase.GetAllAdmins(c.UserContext())

	if err != nil {
		return helper.RespondError(c, err)
	}
	if len(admins) == 0 {
		return helper.RespondSuccess(c, fiber.StatusOK, admins, "Admins not found")
	}

	return helper.RespondSuccess(c, fiber.StatusOK, admins, "")
}

// GetAdminByID godoc
// @Summary      Get admin user
// @Description  ดึงข้อมูลผู้ดูแลระบบตาม ID
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Admin ID"
// @Success      200  {object}  common.BaseResponse[entity.AdminRes]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      403  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/admin/{id} [get]
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

	admin, err := h.usecase.GetAdminByID(c.UserContext(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, admin, "")
}

// CreateAdmin godoc
// @Summary      Create admin user
// @Description  สร้างผู้ดูแลระบบใหม่
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      entity.AdminReq  true  "Admin data"
// @Success      201  {object}  common.BaseResponse[any]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      403  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/admin [post]
func (h *AdminController) CreateAdmin(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	admin := &entity.AdminReq{}
	if err := c.BodyParser(admin); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse request body")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	creatorID := userID

	adminID, err := h.usecase.CreateAdmin(c.UserContext(), admin, &creatorID)
	if err != nil {
		h.logger.Error().Err(err).Str("username", admin.Username).Msg("Failed to create admin in usecase")
		return helper.RespondError(c, err)
	}
	return helper.RespondSuccess(c, fiber.StatusCreated, common.IDResponse{ID: adminID}, "Admin created successfully")

}

// UpdateAdmin godoc
// @Summary      Update admin user
// @Description  อัปเดตข้อมูลผู้ดูแลระบบ
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int             true  "Admin ID"
// @Param        body  body      entity.AdminReq  true  "Admin data"
// @Success      200   {object}  common.BaseResponse[any]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      403   {object}  common.BaseErrorResponse
// @Failure      404   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/admin/{id} [put]
func (h *AdminController) UpdateAdmin(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	idStr := c.Params("id")
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid admin ID format")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"id": "invalid admin ID format"}))
	}

	adminReq := new(entity.AdminReq)
	if err := c.BodyParser(adminReq); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse request body")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	if err := h.usecase.UpdateAdmin(c.UserContext(), adminID, userID, adminReq); err != nil {
		h.logger.Error().Err(err).Int("admin_id", adminID).Msg("Failed to update admin")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("admin_id", adminID).Msg("Admin updated successfully")
	return helper.RespondSuccess[any](c, fiber.StatusOK, nil, "Admin updated successfully")
}

// DeleteAdmin godoc
// @Summary      Delete admin user
// @Description  ลบผู้ดูแลระบบ
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Admin ID"
// @Success      200  {object}  common.BaseResponse[any]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      403  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/admin/{id} [delete]
func (h *AdminController) DeleteAdmin(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	idStr := c.Params("id")
	adminID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid admin ID format")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"id": "invalid admin ID format"}))
	}

	if err := h.usecase.DeleteAdmin(c.UserContext(), adminID, userID); err != nil {
		h.logger.Error().Err(err).Int("admin_id", adminID).Msg("Failed to delete admin")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("admin_id", adminID).Msg("Admin deleted successfully")
	return helper.RespondSuccess[any](c, fiber.StatusOK, nil, "Admin deleted successfully")
}
