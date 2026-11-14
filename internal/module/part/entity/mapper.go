package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbPartToPartRes(dbPart db.GetPartByIDRow) *PartRes {
	return &PartRes{
		PartID:      int32(dbPart.PartID),
		PartBrandID: int32(dbPart.PartBrandID),
		PartTypeID:  int32(dbPart.PartTypeID),
		Name:        dbPart.Name,
		SKU:         dbPart.Sku,
		OemNumber:   utils.NullStringToStringPtr(dbPart.OemNumber),
		Description: utils.NullStringToStringPtr(dbPart.Description),
		Price:       int32(dbPart.Price),
		Quantity:    int32(dbPart.Quantity),
		WeightG:     utils.NullInt32ToInt32Ptr(dbPart.WeightG),
	}
}

func MapDbPartBySKUToPartRes(dbPart db.GetPartBySKURow) *PartRes {
	return &PartRes{
		PartID:      int32(dbPart.PartID),
		PartBrandID: int32(dbPart.PartBrandID),
		PartTypeID:  int32(dbPart.PartTypeID),
		Name:        dbPart.Name,
		SKU:         dbPart.Sku,
		OemNumber:   utils.NullStringToStringPtr(dbPart.OemNumber),
		Description: utils.NullStringToStringPtr(dbPart.Description),
		Price:       int32(dbPart.Price),
		Quantity:    int32(dbPart.Quantity),
		WeightG:     utils.NullInt32ToInt32Ptr(dbPart.WeightG),
	}
}

func MapPartToPartParam(req *PartReq, createdBy *int) db.CreatePartParams {
	return db.CreatePartParams{
		PartBrandID:      int32(req.PartBrandID),
		PartTypeID:       int32(req.PartTypeID),
		Name:             req.Name,
		Sku:              req.SKU,
		OemNumber:        utils.StringPtrToNullString(req.OemNumber),
		Description:      utils.StringPtrToNullString(req.Description),
		Price:            int32(req.Price),
		Quantity:         int32(req.Quantity),
		CreatedByAdminID: utils.IntNullToNullInt32(createdBy),
		CreatedAt:        utils.NowUTC(),
		WeightG:          utils.Int32PtrToNullInt32(req.WeightG),
	}
}

func MapDbPartsToPartRes(dbPart db.ListPartsRow) *PartRes {
	return &PartRes{
		PartID:      int32(dbPart.PartID),
		PartBrandID: int32(dbPart.PartBrandID),
		PartTypeID:  int32(dbPart.PartTypeID),
		Name:        dbPart.Name,
		SKU:         dbPart.Sku,
		OemNumber:   utils.NullStringToStringPtr(dbPart.OemNumber),
		Description: utils.NullStringToStringPtr(dbPart.Description),
		Price:       int32(dbPart.Price),
		Quantity:    int32(dbPart.Quantity),
		WeightG:     utils.NullInt32ToInt32Ptr(dbPart.WeightG),
	}
}

func MapUpdatePartParams(id int, req *PartReq, updatedBy int) db.UpdatePartByIDParams {
	return db.UpdatePartByIDParams{
		PartID:           int32(id),
		PartBrandID:      int32(req.PartBrandID),
		PartTypeID:       int32(req.PartTypeID),
		Name:             req.Name,
		Sku:              req.SKU,
		OemNumber:        utils.StringPtrToNullString(req.OemNumber),
		Description:      utils.StringPtrToNullString(req.Description),
		Price:            int32(req.Price),
		Quantity:         int32(req.Quantity),
		IsActive:         req.IsActive,
		UpdatedByAdminID: utils.IntNullToNullInt32(&updatedBy),
		UpdatedAt:        utils.NullTimeNow(),
		WeightG:          utils.Int32PtrToNullInt32(req.WeightG),
	}
}

func MapUpdatePartIsActiveParams(id int, isActive bool, updatedBy int) db.DeletePartByIDParams {
	return db.DeletePartByIDParams{
		UpdatedByAdminID: utils.IntNullToNullInt32(&updatedBy),
		UpdatedAt:        utils.NullTimeNow(),
		PartID:           int32(id),
	}
}
