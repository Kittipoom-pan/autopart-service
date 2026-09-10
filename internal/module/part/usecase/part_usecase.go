package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
)

type PartUsecase interface {
	GetPartByID(ctx context.Context, id int) (*entity.PartRes, error)
	GetAllParts(ctx context.Context) ([]*entity.PartRes, error)
	CreatePart(ctx context.Context, user *entity.PartReq, userID *int) (int64, error)
	UpdatePart(ctx context.Context, id int, user *entity.PartReq, userID int) error
	DeletePart(ctx context.Context, id int, userID int) error
}
