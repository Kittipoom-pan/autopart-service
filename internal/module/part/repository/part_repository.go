package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
)

type PartRepository interface {
	GetPartByID(ctx context.Context, id int) (*entity.PartRes, error)
	GetPartBySKU(ctx context.Context, sku string) (*entity.PartRes, error)
	CreatePart(ctx context.Context, part *entity.PartReq, createdBy *int) (int64, error)
	GetAllParts(ctx context.Context) ([]*entity.PartRes, error)
	UpdatePart(ctx context.Context, id int, part *entity.PartReq, updatedBy int) error
	DeletePart(ctx context.Context, id int, updatedBy int) error
	GetPartsByIDs(ctx context.Context, ids []int32) (map[int32]*entity.PartRes, error)
}
