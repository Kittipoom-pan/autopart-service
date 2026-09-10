package repository

import (
	"context"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entitie"
)

type PartRepository interface {
	GetPartByID(ctx context.Context, id int) (*entitie.PartRes, error)
	CreatePart(ctx context.Context, param db.CreatePartParams) (int64, error)
	GetAllParts(ctx context.Context) ([]*entitie.PartRes, error)
	UpdatePart(ctx context.Context, param db.UpdatePartByIDParams) error
	DeletePart(ctx context.Context, param db.DeletePartByIDParams) error
}
