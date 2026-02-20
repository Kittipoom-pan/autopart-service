package entitie

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
		Description: utils.ToStringPtr(dbPart.Description),
		Price:       utils.ToInt32Ptr(dbPart.Price),
		Quantity:    utils.ToInt32Ptr(dbPart.Quantity),
	}
}

func MapPartToPartParam(req *PartReq, createdBy string) db.CreatePartParams {
	return db.CreatePartParams{
		PartBrandID: int32(req.PartBrandID),
		PartTypeID:  int32(req.PartTypeID),
		Name:        req.Name,
		Sku:         req.SKU,
		Description: utils.StringPtrToNullString(req.Description),
		Price:       utils.IntToNullInt32(req.Price),
		Quantity:    utils.IntToNullInt32(req.Quantity),
		CreatedBy:   utils.StringToNullString(createdBy),
		CreatedAt:   utils.NullTimeNow(),
	}
}

func MapDbPartsToPartRes(dbPart db.ListPartsRow) *PartRes {
	return &PartRes{
		PartID:      uint32(dbPart.PartID),
		PartBrandID: uint32(dbPart.PartBrandID),
		PartTypeID:  uint32(dbPart.PartTypeID),
		Name:        dbPart.Name,
		SKU:         dbPart.Sku,
		Description: utils.ToStringPtr(dbPart.Description),
		Price:       utils.ToInt32Ptr(dbPart.Price),
		Quantity:    utils.ToInt32Ptr(dbPart.Quantity),
	}
}

func MapUpdatePartParams(id int, req *PartReq, updatedBy string) db.UpdatePartByIDParams {
	return db.UpdatePartByIDParams{
		PartID:      int32(id),
		PartBrandID: int32(req.PartBrandID),
		PartTypeID:  int32(req.PartTypeID),
		Name:        req.Name,
		Sku:         req.SKU,
		Description: utils.StringPtrToNullString(req.Description),
		Price:       utils.IntToNullInt32(req.Price),
		Quantity:    utils.IntToNullInt32(req.Quantity),
		IsActive:    req.IsActive,
		UpdatedBy:   utils.StringToNullString(updatedBy),
		UpdatedAt:   utils.NullTimeNow(),
	}
}

func MapUpdatePartIsActiveParams(id int, isActive bool, updatedBy string) db.DeletePartByIDParams {
	return db.DeletePartByIDParams{
		UpdatedBy: utils.StringToNullString(updatedBy),
		UpdatedAt: utils.NullTimeNow(),
		PartID:    int32(id),
	}
}
