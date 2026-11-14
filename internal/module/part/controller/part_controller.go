package controller

import (
	"encoding/json"
	"mime/multipart"
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/usecase"
	parterror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type PartController struct {
	usecase usecase.PartUsecase
	logger  zerolog.Logger
}

func NewPartController(usecase usecase.PartUsecase, logger zerolog.Logger) *PartController {
	return &PartController{
		usecase: usecase,
		logger:  logger,
	}
}

// GetAllParts godoc
// @Summary      List parts
// @Description  ดึงรายการอะไหล่ทั้งหมด
// @Tags         part
// @Produce      json
// @Success      200  {object}  common.BaseResponse[[]entity.PartRes]
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/part [get]
func (h *PartController) GetAllParts(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get all parts request")

	parts, err := h.usecase.GetAllParts(c.UserContext())
	if err != nil {
		return helper.RespondError(c, err)
	}

	if parts == nil {
		parts = []*entity.PartRes{}
	}

	return helper.RespondSuccess(c, fiber.StatusOK, parts, "Parts retrieved successfully")
}

// GetPartByID godoc
// @Summary  Get part
// @Tags     part
// @Produce  json
// @Param    id   path      int  true  "Part ID"
// @Success 200 {object} common.BaseResponse[entity.PartRes]
// @Failure  400  {object}  common.BaseErrorResponse
// @Failure  404  {object}  common.BaseErrorResponse
// @Failure  500  {object}  common.BaseErrorResponse
// @Router   /v1/part/{id} [get]
func (h *PartController) GetPartByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("part_id", idStr).Msg("Get part request")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Invalid part id format")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"id": "invalid part id format"}))
	}

	part, err := h.usecase.GetPartByID(c.UserContext(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, part, "")
}

func (h *PartController) GetPartBySKU(c *fiber.Ctx) error {
	sku := c.Params("sku")
	h.logger.Debug().Str("sku", sku).Msg("Get part by SKU request")

	if sku == "" {
		h.logger.Warn().Msg("SKU parameter is missing")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"sku": "sku is required"}))
	}

	part, err := h.usecase.GetPartBySKU(c.UserContext(), sku)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, part, "")
}

// CreatePart godoc
// @Summary      Create part
// @Description  สร้างอะไหล่ใหม่ พร้อมอัปโหลดรูปได้ ต้องใช้ JWT ส่งข้อมูล part เป็น JSON string ในฟิลด์ "part" และไฟล์รูปในฟิลด์ "images" (multipart/form-data)
// @Tags         part
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        part    formData  string  true   "Part data as JSON string"
// @Param        images  formData  []file  false  "Part images (.jpg, .jpeg, .png, .webp) max 5MB each" collectionFormat(multi)
// @Success 201 {object} common.BaseResponse[common.IDResponse]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      403   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/part [post]
func (h *PartController) CreatePart(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	partReq, err := parsePartForm(c)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse part data")
		return helper.RespondError(c, err)
	}

	images, err := parsePartImages(c)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse part images")
		return helper.RespondError(c, err)
	}

	partID, err := h.usecase.CreatePart(c.UserContext(), partReq, &userID, images)
	if err != nil {
		h.logger.Error().Err(err).Str("sku", partReq.SKU).Msg("Failed to create part")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int64("part_id", partID).Msg("Part created successfully")
	return helper.RespondSuccess(c, fiber.StatusCreated, common.IDResponse{ID: partID}, "Part created successfully")
}

// UpdatePart godoc
// @Summary      Update part
// @Description  อัปเดตอะไหล่ พร้อมอัปโหลดรูปได้ ต้องใช้ JWT ส่งข้อมูล part เป็น JSON string ในฟิลด์ "part" และไฟล์รูปในฟิลด์ "images" (multipart/form-data) หากส่งรูปมาจะแทนที่รูปเดิมทั้งหมด
// @Tags         part
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int     true   "Part ID"
// @Param        part    formData  string  true   "Part data as JSON string"
// @Param        images  formData  []file  false  "Part images (.jpg, .jpeg, .png, .webp) max 5MB each" collectionFormat(multi)
// @Success      200   {object}  common.BaseResponse[any]
// @Failure      400   {object}  common.BaseErrorResponse
// @Failure      401   {object}  common.BaseErrorResponse
// @Failure      403   {object}  common.BaseErrorResponse
// @Failure      404   {object}  common.BaseErrorResponse
// @Failure      500   {object}  common.BaseErrorResponse
// @Router       /v1/part/{id} [put]
func (h *PartController) UpdatePart(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid part ID format")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"id": "invalid part ID format"}))
	}

	partReq, err := parsePartForm(c)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse part data")
		return helper.RespondError(c, err)
	}

	images, err := parsePartImages(c)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse part images")
		return helper.RespondError(c, err)
	}

	if err := h.usecase.UpdatePart(c.UserContext(), id, partReq, userID, images); err != nil {
		h.logger.Error().Err(err).Int("part_id", id).Msg("Failed to update part")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("part_id", id).Msg("Part updated successfully")
	return helper.RespondSuccess[any](c, fiber.StatusOK, nil, "Part updated successfully")
}

// DeletePart godoc
// @Summary      Delete part
// @Tags         part
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Part ID"
// @Success      200  {object}  common.BaseResponse[any]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      403  {object}  common.BaseErrorResponse
// @Failure      404  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/part/{id} [delete]
func (h *PartController) DeletePart(c *fiber.Ctx) error {
	userID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid part ID format")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"id": "invalid part ID format"}))
	}

	if err := h.usecase.DeletePart(c.UserContext(), id, userID); err != nil {
		h.logger.Error().Err(err).Int("part_id", id).Msg("Failed to delete part")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("part_id", id).Msg("Part deleted successfully")
	return helper.RespondSuccess[any](c, fiber.StatusOK, nil, "Part deleted successfully")
}

func parsePartForm(c *fiber.Ctx) (*entity.PartReq, error) {
	raw := c.FormValue("part")
	if raw == "" {
		return nil, parterror.InvalidRequestData(map[string]string{"part": "part data is required"})
	}

	partReq := &entity.PartReq{}
	if err := json.Unmarshal([]byte(raw), partReq); err != nil {
		return nil, parterror.InvalidRequestData(map[string]string{"part": "invalid part data"})
	}

	return partReq, nil
}

func parsePartImages(c *fiber.Ctx) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, parterror.InvalidRequestData(map[string]string{"images": "invalid multipart form data"})
	}

	return form.File["images"], nil
}
