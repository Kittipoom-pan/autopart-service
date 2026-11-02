package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func parseBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case []byte:
		return len(val) > 0 && val[0] == 1
	case int64:
		return val == 1
	case int32:
		return val == 1
	case int:
		return val == 1
	default:
		return false
	}
}

func MapDbAddressToAddressRes(dbAddr db.Address) *AddressRes {
	return &AddressRes{
		AddressID:     dbAddr.AddressID,
		CustomerID:    dbAddr.CustomerID,
		Title:         utils.NullStringToStringPtr(dbAddr.Title),
		RecipientName: dbAddr.RecipientName,
		AddressLine:   dbAddr.AddressLine,
		SubdistrictID: dbAddr.SubdistrictID,
		ZipCode:       dbAddr.ZipCode,
		PhoneNumber:   dbAddr.PhoneNumber,
		IsDefault:     parseBool(dbAddr.IsDefault),
		CreatedAt:     dbAddr.CreatedAt,
		UpdatedAt:     utils.NullTimeToTimePtr(dbAddr.UpdatedAt),
	}
}

func MapAddressDetailToAddressRes(row db.GetAddressByIDAndCustomerRow) *AddressRes {
	provinceID := row.ProvinceID
	return &AddressRes{
		AddressID:     row.AddressID,
		CustomerID:    row.CustomerID,
		Title:         utils.NullStringToStringPtr(row.Title),
		RecipientName: row.RecipientName,
		AddressLine:   row.AddressLine,
		SubdistrictID: row.SubdistrictID,
		Subdistrict:   row.Subdistrict,
		District:      row.District,
		ProvinceID:    &provinceID,
		Province:      row.Province,
		ZipCode:       row.ZipCode,
		PhoneNumber:   row.PhoneNumber,
		IsDefault:     parseBool(row.IsDefault),
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     utils.NullTimeToTimePtr(row.UpdatedAt),
	}
}

func MapDbAddressesToAddressResList(dbAddrs []db.Address) []*AddressRes {
	items := make([]*AddressRes, 0, len(dbAddrs))
	for _, a := range dbAddrs {
		items = append(items, MapDbAddressToAddressRes(a))
	}
	return items
}

func MapCreateAddressReqToParams(req *CreateAddressReq, customerID int32, createdBy *int32) db.CreateAddressParams {
	return db.CreateAddressParams{
		CustomerID:    customerID,
		Title:         utils.StringPtrToNullString(req.Title),
		RecipientName: req.RecipientName,
		AddressLine:   req.AddressLine,
		SubdistrictID: req.SubdistrictID,
		ZipCode:       req.ZipCode,
		PhoneNumber:   req.PhoneNumber,
		IsDefault:     req.IsDefault,
		CreatedAt:     utils.NowUTC(),
		CreatedBy:     utils.Int32PtrToNullInt32(createdBy),
	}
}

func MapUpdateAddressReqToParams(req *UpdateAddressReq, addressID int32, customerID int32, updatedBy *int32) db.UpdateAddressParams {
	return db.UpdateAddressParams{
		Title:         utils.StringPtrToNullString(req.Title),
		RecipientName: req.RecipientName,
		AddressLine:   req.AddressLine,
		SubdistrictID: req.SubdistrictID,
		ZipCode:       req.ZipCode,
		PhoneNumber:   req.PhoneNumber,
		IsDefault:     false,
		UpdatedAt:     utils.NullTimeNow(),
		UpdatedBy:     utils.Int32PtrToNullInt32(updatedBy),
		AddressID:     addressID,
		CustomerID:    customerID,
	}
}
