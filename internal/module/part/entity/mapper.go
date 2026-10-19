package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbPartToPartRes(dbPart db.GetPartByIDRow) *PartRes {
	return &PartRes{
		PartID:      uint32(dbPart.PartID),
		PartBrandID: uint32(dbPart.PartBrandID),
		PartTypeID:  uint32(dbPart.PartTypeID),
		Name:        dbPart.Name,
		SKU:         dbPart.Sku,
		Description: utils.NullStringToStringPtr(dbPart.Description),
		Price:       int(dbPart.Price),
		Quantity:    int(dbPart.Quantity),
	}
}

func MapDbPartBySKUToPartRes(dbPart db.GetPartBySKURow) *PartRes {
	return &PartRes{
		PartID:      uint32(dbPart.PartID),
		PartBrandID: uint32(dbPart.PartBrandID),
		PartTypeID:  uint32(dbPart.PartTypeID),
		Name:        dbPart.Name,
		SKU:         dbPart.Sku,
		Description: utils.NullStringToStringPtr(dbPart.Description),
		Price:       int(dbPart.Price),
		Quantity:    int(dbPart.Quantity),
	}
}

func MapPartToPartParam(req *PartReq, createdBy *int) db.CreatePartParams {
	return db.CreatePartParams{
		PartBrandID: int32(req.PartBrandID),
		PartTypeID:  int32(req.PartTypeID),
		Name:        req.Name,
		Sku:         req.SKU,
		Description: utils.StringPtrToNullString(req.Description),
		Price:       int32(req.Price),
		Quantity:    int32(req.Quantity),
		CreatedBy:   utils.IntNullToNullInt32(createdBy),
		CreatedAt:   utils.NowUTC(),
	}
}

func MapDbPartsToPartRes(dbPart db.ListPartsRow) *PartRes {
	return &PartRes{
		PartID:      uint32(dbPart.PartID),
		PartBrandID: uint32(dbPart.PartBrandID),
		PartTypeID:  uint32(dbPart.PartTypeID),
		Name:        dbPart.Name,
		SKU:         dbPart.Sku,
		Description: utils.NullStringToStringPtr(dbPart.Description),
		Price:       int(dbPart.Price),
		Quantity:    int(dbPart.Quantity),
	}
}

func MapUpdatePartParams(id int, req *PartReq, updatedBy int) db.UpdatePartByIDParams {
	return db.UpdatePartByIDParams{
		PartID:      int32(id),
		PartBrandID: int32(req.PartBrandID),
		PartTypeID:  int32(req.PartTypeID),
		Name:        req.Name,
		Sku:         req.SKU,
		Description: utils.StringPtrToNullString(req.Description),
		Price:       int32(req.Price),
		Quantity:    int32(req.Quantity),
		IsActive:    req.IsActive,
		UpdatedBy:   utils.IntNullToNullInt32(&updatedBy),
		UpdatedAt:   utils.NullTimeNow(),
	}
}

func MapUpdatePartIsActiveParams(id int, isActive bool, updatedBy int) db.DeletePartByIDParams {
	return db.DeletePartByIDParams{
		UpdatedBy: utils.IntNullToNullInt32(&updatedBy),
		UpdatedAt: utils.NullTimeNow(),
		PartID:    int32(id),
	}
}
