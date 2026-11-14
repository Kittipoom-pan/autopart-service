package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbImageToImage(img db.Image) Image {
	return Image{
		ImageID:       img.ImageID,
		ReferenceID:   img.ReferenceID,
		ReferenceType: img.ReferenceType,
		FileName:      img.FileName,
		ImageURL:      img.ImageUrl,
		AltText:       img.AltText,
		IsPrimary:     utils.ParseBool(img.IsPrimary),
		SortImage:     img.SortImage,
		CreatedAt:     img.CreatedAt,
	}
}

func MapImageToImageRes(img Image) ImageRes {
	return ImageRes{
		ImageID:       img.ImageID,
		ReferenceID:   img.ReferenceID,
		ReferenceType: img.ReferenceType,
		FileName:      img.FileName,
		ImageURL:      img.ImageURL,
		AltText:       img.AltText,
		IsPrimary:     img.IsPrimary,
		SortImage:     img.SortImage,
		CreatedAt:     img.CreatedAt,
	}
}
