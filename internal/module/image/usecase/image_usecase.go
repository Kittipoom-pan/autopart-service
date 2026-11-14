package usecase

import (
	"context"
	"mime/multipart"

	"github.com/Kittipoom-pan/autopart-service/internal/module/image/entity"
)

// UploadResult carries the persisted images together with the file names that
// were written to disk, so callers running inside an outer transaction can
// clean up correctly:
//   - SavedFiles    : new files; remove them if the outer transaction rolls back.
//   - ReplacedFiles : old files replaced by this call; remove them only after the
//     outer transaction commits.
type UploadResult struct {
	Images        []entity.ImageRes
	SavedFiles    []string
	ReplacedFiles []string
}

type ImageUsecase interface {
	UploadImage(ctx context.Context, adminID int, file *multipart.FileHeader, req *entity.UploadImageReq) (*entity.ImageRes, error)
	UploadImages(ctx context.Context, adminID int, files []*multipart.FileHeader, req *entity.UploadImageReq) (*UploadResult, error)
	ReplaceImages(ctx context.Context, adminID int, files []*multipart.FileHeader, req *entity.UploadImageReq) (*UploadResult, error)
	RemoveFiles(referenceType string, fileNames ...string) error
}
