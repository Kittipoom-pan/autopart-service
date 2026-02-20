package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entitie"
)

type PartRepository interface {
	GetPartByID(ctx context.Context, id int) (*entitie.PartRes, error)
	CreatePart(ctx context.Context, part *entitie.PartReq) (int64, error)
	GetAllParts(ctx context.Context) ([]*entitie.PartRes, error)
	UpdatePart(ctx context.Context, id int, part *entitie.PartReq) error
	DeletePart(ctx context.Context, id int) error
}
