package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/image/entity"
)

type CreateImageInput struct {
	ReferenceID      int32
	ReferenceType    string
	FileName         string
	ImageURL         string
	AltText          string
	IsPrimary        bool
	SortImage        int32
	CreatedByAdminID *int
}

type ImageRepository interface {
	ListImagesByReference(ctx context.Context, referenceType string, referenceID int32) ([]entity.Image, error)
	CreateImage(ctx context.Context, input CreateImageInput) (int64, error)
	DeleteImage(ctx context.Context, imageID int64) error
}
