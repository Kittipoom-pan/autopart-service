package controller

import (
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	_ "github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/usecase"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type ImageController struct {
	usecase usecase.ImageUsecase
	logger  zerolog.Logger
}

func NewImageController(usecase usecase.ImageUsecase, logger zerolog.Logger) *ImageController {
	return &ImageController{
		usecase: usecase,
		logger:  logger,
	}
}

// UploadImage godoc
// @Summary      Upload image
// @Description  อัปโหลดรูปภาพ (multipart/form-data) รองรับ .jpg, .jpeg, .png, .webp ขนาดไม่เกิน 5MB ระบบจะย่อความกว้างไม่เกิน 1200px และแปลงเป็น .webp
// @Tags         image
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file            formData  file    true  "Image file (.jpg, .jpeg, .png, .webp) max 5MB"
// @Param        reference_id    formData  int     true  "Reference ID"
// @Param        reference_type  formData  string  true  "Reference type (part, customer, order, car_brand, car_model, part_brand, part_type)"
// @Param        alt_text        formData  string  true  "Alternative text"
// @Success      201  {object}  common.BaseResponse[entity.ImageRes]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      403  {object}  common.BaseErrorResponse
// @Failure      422  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/image [post]
func (h *ImageController) UploadImage(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Upload image request")

	adminID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	var req entity.UploadImageReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid multipart form data")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid form data"}))
	}

	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Warn().Err(err).Msg("Image file is missing")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"file": "file is required"}))
	}

	res, err := h.usecase.UploadImage(c.UserContext(), adminID, file, &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, res, "Image uploaded successfully")
}

// UploadImages godoc
// @Summary      Upload multiple images
// @Description  อัปโหลดรูปภาพหลายรูปในคำขอเดียว (multipart/form-data) ใช้ field "files" ซ้ำได้ รองรับ .jpg, .jpeg, .png, .webp รูปละไม่เกิน 5MB ระบบจะย่อความกว้างไม่เกิน 1200px และแปลงเป็น .webp
// @Tags         image
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        files           formData  []file  true  "Image files (.jpg, .jpeg, .png, .webp) max 5MB each" collectionFormat(multi)
// @Param        reference_id    formData  int     true  "Reference ID"
// @Param        reference_type  formData  string  true  "Reference type (part, customer, order, car_brand, car_model, part_brand, part_type)"
// @Param        alt_text        formData  string  true  "Alternative text"
// @Success      201  {object}  common.BaseResponse[[]entity.ImageRes]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      403  {object}  common.BaseErrorResponse
// @Failure      413  {object}  common.BaseErrorResponse
// @Failure      422  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/image/multiple [post]
func (h *ImageController) UploadImages(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Upload multiple images request")

	adminID, err := auth.GetUserID(c, h.logger)
	if err != nil {
		return helper.RespondError(c, err)
	}

	var req entity.UploadImageReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Invalid multipart form data")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid form data"}))
	}

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Warn().Err(err).Msg("Invalid multipart form data")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"body": "invalid form data"}))
	}

	files := form.File["files"]
	if len(files) == 0 {
		h.logger.Warn().Msg("Image files are missing")
		return helper.RespondError(c, customerror.InvalidRequestData(map[string]string{"files": "at least one file is required"}))
	}

	results, err := h.usecase.UploadImages(c.UserContext(), adminID, files, &req)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, results.Images, "Images uploaded successfully")
}
