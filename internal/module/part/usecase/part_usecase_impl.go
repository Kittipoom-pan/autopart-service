package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	imageentity "github.com/Kittipoom-pan/autopart-service/internal/module/image/entity"
	imageusecase "github.com/Kittipoom-pan/autopart-service/internal/module/image/usecase"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/usecase/validation"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
)

const (
	skuTypeMaxLen  = 8
	skuBrandMaxLen = 8
	altTextMaxLen  = 150
)

type partUsecase struct {
	repo         repository.PartRepository
	imageUsecase imageusecase.ImageUsecase
	txManager    database.TxManager
	logger       zerolog.Logger
}

func NewPartUsecase(repo repository.PartRepository, imageUsecase imageusecase.ImageUsecase, txManager database.TxManager, logger zerolog.Logger) PartUsecase {
	return &partUsecase{
		repo:         repo,
		imageUsecase: imageUsecase,
		txManager:    txManager,
		logger:       logger,
	}
}

func (u *partUsecase) GetPartByID(ctx context.Context, id int) (*entity.PartRes, error) {
	u.logger.Debug().Int("part_id", id).Msg("GetPartByID")
	part, err := u.repo.GetPartByID(ctx, id)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Part")
	}
	return part, nil
}

func (u *partUsecase) GetPartBySKU(ctx context.Context, sku string) (*entity.PartRes, error) {
	u.logger.Debug().Str("sku", sku).Msg("GetPartBySKU")
	part, err := u.repo.GetPartBySKU(ctx, sku)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Part")
	}
	return part, nil
}

func (u *partUsecase) GetAllParts(ctx context.Context) ([]*entity.PartRes, error) {
	u.logger.Debug().Msg("GetAllParts")
	parts, err := u.repo.GetAllParts(ctx)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Part")
	}
	return parts, nil
}

func (u *partUsecase) CreatePart(ctx context.Context, part *entity.PartReq, userID *int, images []*multipart.FileHeader) (int64, error) {
	if err := validation.ValidatePartRequest(part, false); err != nil {
		u.logger.Warn().Err(err).Msg("part create validation failed")
		return 0, err
	}

	adminID := 0
	if userID != nil {
		adminID = *userID
	}

	var partID int64
	var upload *imageusecase.UploadResult

	err := u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		partTypeCode, err := u.repo.GetPartTypeCode(ctx, part.PartTypeID)
		if err != nil {
			return err
		}
		partBrandName, err := u.repo.GetPartBrandName(ctx, part.PartBrandID)
		if err != nil {
			return err
		}

		sku, err := u.generateSKU(ctx, partTypeCode, partBrandName)
		if err != nil {
			return err
		}
		part.SKU = sku

		id, err := u.repo.CreatePart(ctx, part, userID)
		if err != nil {
			return err
		}
		partID = id

		if len(images) > 0 {
			upload, err = u.uploadPartImages(ctx, adminID, int32(id), sku, part.Name, part.CarBrandID, partBrandName, images, false)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		u.cleanupSavedFiles(upload)
		return 0, helper.MapDBErrorToAPIError(err, "Part")
	}

	u.logger.Info().Int64("part_id", partID).Str("sku", part.SKU).Msg("part created successfully")
	return partID, nil
}

func (u *partUsecase) UpdatePart(ctx context.Context, id int, partReq *entity.PartReq, userID int, images []*multipart.FileHeader) error {
	if err := validation.ValidatePartRequest(partReq, true); err != nil {
		u.logger.Warn().Err(err).Msg("part update validation failed")
		return err
	}

	var upload *imageusecase.UploadResult

	err := u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		existing, err := u.repo.GetPartByID(ctx, id)
		if err != nil {
			return err
		}
		// SKU อ้างอิงตามเดิม ไม่สร้างใหม่ตอน update
		partReq.SKU = existing.SKU

		if err := u.repo.UpdatePart(ctx, id, partReq, userID); err != nil {
			return err
		}

		if len(images) > 0 {
			partBrandName, err := u.repo.GetPartBrandName(ctx, partReq.PartBrandID)
			if err != nil {
				return err
			}
			upload, err = u.uploadPartImages(ctx, userID, int32(id), existing.SKU, partReq.Name, partReq.CarBrandID, partBrandName, images, true)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		u.cleanupSavedFiles(upload)
		return helper.MapDBErrorToAPIError(err, "Part")
	}

	// ลบไฟล์รูปเดิมหลัง commit สำเร็จ
	if upload != nil && len(upload.ReplacedFiles) > 0 {
		if err := u.imageUsecase.RemoveFiles(string(imageentity.ReferenceTypePart), upload.ReplacedFiles...); err != nil {
			u.logger.Warn().Err(err).Msg("failed to remove replaced part image files")
		}
	}

	u.logger.Info().Int("part_id", id).Msg("part updated successfully")
	return nil
}

func (u *partUsecase) DeletePart(ctx context.Context, id int, userID int) error {
	if err := u.repo.DeletePart(ctx, id, userID); err != nil {
		return helper.MapDBErrorToAPIError(err, "Part")
	}

	u.logger.Info().Int("part_id", id).Msg("part deleted successfully")
	return nil
}

// uploadPartImages สร้าง alt_text และอัปโหลดรูปภายใต้ transaction ปัจจุบัน
// replace = true สำหรับ update (แทนที่รูปเดิมทั้งหมด)
func (u *partUsecase) uploadPartImages(
	ctx context.Context,
	adminID int,
	referenceID int32,
	filePrefix, name string,
	carBrandID int32,
	partBrandName string,
	images []*multipart.FileHeader,
	replace bool,
) (*imageusecase.UploadResult, error) {
	carBrandName := ""
	if carBrandID > 0 {
		brandName, err := u.repo.GetCarBrandName(ctx, carBrandID)
		if err != nil {
			return nil, err
		}
		carBrandName = brandName
	}

	req := &imageentity.UploadImageReq{
		ReferenceID:   referenceID,
		ReferenceType: string(imageentity.ReferenceTypePart),
		AltText:       buildAltText(name, carBrandName, partBrandName),
		FilePrefix:    filePrefix,
	}

	if replace {
		return u.imageUsecase.ReplaceImages(ctx, adminID, images, req)
	}
	return u.imageUsecase.UploadImages(ctx, adminID, images, req)
}

// generateSKU สร้าง SKU รูปแบบ [TypePrefix]-[BrandPrefix]-[RunningNo]
func (u *partUsecase) generateSKU(ctx context.Context, typeCode, partBrandName string) (string, error) {
	typePrefix := normalizeSKUPart(typeCode, skuTypeMaxLen)
	brandPrefix := normalizeSKUPart(partBrandName, skuBrandMaxLen)
	if typePrefix == "" || brandPrefix == "" {
		return "", customerror.InvalidRequestData(map[string]string{
			"sku": "cannot generate SKU from part type and part brand",
		})
	}

	prefix := typePrefix + "-" + brandPrefix // ยาวสุด 17 ตัว ไม่ต้อง truncate อีก

	nextNumber, err := u.repo.NextSKUNumber(ctx, prefix)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%05d", prefix, nextNumber), nil
}

func (u *partUsecase) cleanupSavedFiles(upload *imageusecase.UploadResult) {
	if upload == nil || len(upload.SavedFiles) == 0 {
		return
	}
	if err := u.imageUsecase.RemoveFiles(string(imageentity.ReferenceTypePart), upload.SavedFiles...); err != nil {
		u.logger.Warn().Err(err).Msg("failed to remove orphaned part image files")
	}
}

// เก็บเฉพาะ A-Z และ 0-9 กันอักขระไทย/สัญลักษณ์หลุดเข้า SKU
func normalizeSKUPart(name string, maxLen int) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(name) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			if b.Len() >= maxLen {
				break
			}
		}
	}
	return b.String()
}

func buildAltText(name, carBrandName, partBrandName string) string {
	parts := make([]string, 0, 3)
	for _, v := range []string{name, carBrandName, partBrandName} {
		if v = strings.TrimSpace(v); v != "" {
			parts = append(parts, v)
		}
	}
	return truncateRunes(strings.Join(parts, " "), altTextMaxLen)
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}
