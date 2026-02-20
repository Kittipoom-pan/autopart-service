package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entitie"
)

type PartUsecase interface {
	GetPartByID(ctx context.Context, id int) (*entitie.PartRes, error)
	GetAllParts(ctx context.Context) ([]*entitie.PartRes, error)
	CreatePart(ctx context.Context, user *entitie.PartReq) (int64, error)
	UpdatePart(ctx context.Context, id int, user *entitie.PartReq) error
	DeletePart(ctx context.Context, id int) error
}
