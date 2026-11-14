package usecase

import (
	"context"
	"mime/multipart"

	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
)

type PartUsecase interface {
	GetPartByID(ctx context.Context, id int) (*entity.PartRes, error)
	GetPartBySKU(ctx context.Context, sku string) (*entity.PartRes, error)
	GetAllParts(ctx context.Context) ([]*entity.PartRes, error)
	CreatePart(ctx context.Context, user *entity.PartReq, userID *int, images []*multipart.FileHeader) (int64, error)
	UpdatePart(ctx context.Context, id int, user *entity.PartReq, userID int, images []*multipart.FileHeader) error
	DeletePart(ctx context.Context, id int, userID int) error
}
