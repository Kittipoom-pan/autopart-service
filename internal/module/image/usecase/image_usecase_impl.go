package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/usecase/validation"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/Kittipoom-pan/autopart-service/pkg/imagestore"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const suffixAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

type imageUsecase struct {
	repo         repository.ImageRepository
	store        imagestore.Store
	txManager    database.TxManager
	maxSizeBytes int64
	maxFiles     int
	logger       zerolog.Logger
}

func NewImageUsecase(
	repo repository.ImageRepository,
	store imagestore.Store,
	txManager database.TxManager,
	maxSizeBytes int64,
	maxFiles int,
	logger zerolog.Logger,
) ImageUsecase {
	return &imageUsecase{
		repo:         repo,
		store:        store,
		txManager:    txManager,
		maxSizeBytes: maxSizeBytes,
		maxFiles:     maxFiles,
		logger:       logger,
	}
}

func (u *imageUsecase) UploadImage(ctx context.Context, adminID int, file *multipart.FileHeader, req *entity.UploadImageReq) (*entity.ImageRes, error) {
	if file == nil {
		return nil, customerror.InvalidRequestData(map[string]string{"file": "file is required"})
	}

	result, err := u.upload(ctx, adminID, []*multipart.FileHeader{file}, req, false)
	if err != nil {
		return nil, err
	}

	return &result.Images[0], nil
}

func (u *imageUsecase) UploadImages(ctx context.Context, adminID int, files []*multipart.FileHeader, req *entity.UploadImageReq) (*UploadResult, error) {
	if err := u.validateFiles(files); err != nil {
		return nil, err
	}

	return u.upload(ctx, adminID, files, req, false)
}

func (u *imageUsecase) ReplaceImages(ctx context.Context, adminID int, files []*multipart.FileHeader, req *entity.UploadImageReq) (*UploadResult, error) {
	if err := u.validateFiles(files); err != nil {
		return nil, err
	}

	return u.upload(ctx, adminID, files, req, true)
}

// RemoveFiles deletes image files from storage.
func (u *imageUsecase) RemoveFiles(referenceType string, fileNames ...string) error {
	return u.store.Remove(referenceType, fileNames...)
}

func (u *imageUsecase) validateFiles(files []*multipart.FileHeader) error {
	if len(files) == 0 {
		return customerror.InvalidRequestData(map[string]string{"files": "at least one file is required"})
	}

	if len(files) > u.maxFiles {
		return customerror.InvalidRequestData(map[string]string{
			"files": fmt.Sprintf("must not exceed %d files per request", u.maxFiles),
		})
	}

	return nil
}

func (u *imageUsecase) upload(ctx context.Context, adminID int, files []*multipart.FileHeader, req *entity.UploadImageReq, replace bool) (*UploadResult, error) {
	u.logger.Debug().
		Str("reference_type", req.ReferenceType).
		Int32("reference_id", req.ReferenceID).
		Int("file_count", len(files)).
		Msg("UploadImages")

	// 1) ตรวจสอบเงื่อนไขทั้งหมดก่อน ยังไม่เขียนไฟล์ลงดิสก์
	for _, file := range files {
		if err := validation.ValidateUploadImageRequest(req, file.Filename, file.Size, u.maxSizeBytes); err != nil {
			u.logger.Warn().Err(err).Msg("upload image validation failed")
			return nil, err
		}
	}

	// 2) อ่านไฟล์ทั้งหมดเข้าหน่วยความจำ
	payloads := make([][]byte, len(files))
	for i, file := range files {
		data, err := readFileWithLimit(file, u.maxSizeBytes)
		if err != nil {
			u.logger.Warn().Err(err).Str("file", file.Filename).Msg("failed to read uploaded file")
			return nil, customerror.InvalidRequestData(map[string]string{"file": "failed to read uploaded file"})
		}
		if int64(len(data)) > u.maxSizeBytes {
			return nil, customerror.InvalidRequestData(map[string]string{
				"file": fmt.Sprintf("file size must not exceed %dMB", u.maxSizeBytes/(1024*1024)),
			})
		}
		payloads[i] = data
	}

	// หากถูกเรียกจาก transaction ภายนอก ให้เจ้าของ transaction เป็นผู้จัดการไฟล์
	// (ลบไฟล์ใหม่เมื่อ rollback และลบไฟล์เดิมหลัง commit)
	nested := isInTransaction(ctx)

	result := &UploadResult{}
	err := u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		existing, err := u.repo.ListImagesByReference(ctx, req.ReferenceType, req.ReferenceID)
		if err != nil {
			return err
		}

		baseSort := int32(len(existing))
		if replace {
			for _, img := range existing {
				if err := u.repo.DeleteImage(ctx, img.ImageID); err != nil {
					return err
				}
				result.ReplacedFiles = append(result.ReplacedFiles, img.FileName)
			}
			baseSort = 0
		}

		for i := range files {
			savedFile, err := u.store.Save(req.ReferenceType, u.buildBaseName(req.FilePrefix), payloads[i])
			if err != nil {
				return err
			}
			result.SavedFiles = append(result.SavedFiles, savedFile.FileName)

			isPrimary := baseSort == 0 && i == 0
			sortImage := baseSort + int32(i)

			imageID, err := u.repo.CreateImage(ctx, repository.CreateImageInput{
				ReferenceID:      req.ReferenceID,
				ReferenceType:    req.ReferenceType,
				FileName:         savedFile.FileName,
				ImageURL:         savedFile.ImageURL,
				AltText:          req.AltText,
				IsPrimary:        isPrimary,
				SortImage:        sortImage,
				CreatedByAdminID: &adminID,
			})
			if err != nil {
				return err
			}

			result.Images = append(result.Images, entity.ImageRes{
				ImageID:       imageID,
				ReferenceID:   req.ReferenceID,
				ReferenceType: req.ReferenceType,
				FileName:      savedFile.FileName,
				ImageURL:      savedFile.ImageURL,
				AltText:       req.AltText,
				IsPrimary:     isPrimary,
				SortImage:     sortImage,
				CreatedAt:     time.Now(),
			})
		}

		return nil
	})
	if err != nil {
		// transaction ล้มเหลว: ลบไฟล์ใหม่ทั้งหมดที่เพิ่งเขียนลงดิสก์
		if len(result.SavedFiles) > 0 {
			_ = u.store.Remove(req.ReferenceType, result.SavedFiles...)
		}
		if errors.Is(err, imagestore.ErrUnsupportedFormat) {
			return nil, customerror.InvalidRequestData(map[string]string{"file": "unsupported or corrupted image file"})
		}
		return nil, helper.MapDBErrorToAPIError(err, "Image")
	}

	// ลบไฟล์รูปเดิมออกจากดิสก์หลัง commit สำเร็จ (เฉพาะกรณี replace ที่ไม่ได้อยู่ใน tx ภายนอก)
	if !nested && len(result.ReplacedFiles) > 0 {
		if err := u.store.Remove(req.ReferenceType, result.ReplacedFiles...); err != nil {
			u.logger.Warn().Err(err).Msg("failed to remove replaced image files")
		}
	}

	u.logger.Info().
		Int("image_count", len(result.Images)).
		Str("reference_type", req.ReferenceType).
		Int32("reference_id", req.ReferenceID).
		Msg("images uploaded successfully")

	return result, nil
}

// buildBaseName สร้างชื่อไฟล์ รูปแบบ [prefix]-[random] หรือ uuid เมื่อไม่มี prefix
func (u *imageUsecase) buildBaseName(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return uuid.NewString()
	}
	return prefix + "-" + randomSuffix(5)
}

func randomSuffix(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return strings.Repeat("0", n)
	}
	for i := range buf {
		buf[i] = suffixAlphabet[int(buf[i])%len(suffixAlphabet)]
	}
	return string(buf)
}

func isInTransaction(ctx context.Context) bool {
	_, ok := database.GetTx(ctx)
	return ok
}

func readFileWithLimit(file *multipart.FileHeader, maxSizeBytes int64) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	return io.ReadAll(io.LimitReader(src, maxSizeBytes+1))
}
