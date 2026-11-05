package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbAddressToAddressRes(dbAddr db.ListAddressesByCustomerIDRow) *AddressRes {
	return &AddressRes{
		AddressID:     dbAddr.AddressID,
		CustomerID:    dbAddr.CustomerID,
		Title:         utils.NullStringToStringPtr(dbAddr.Title),
		RecipientName: dbAddr.RecipientName,
		AddressLine:   dbAddr.AddressLine,
		ZipCode:       dbAddr.ZipCode,
		PhoneNumber:   dbAddr.PhoneNumber,
		IsDefault:     utils.ParseBool(dbAddr.IsDefault),
		CreatedAt:     dbAddr.CreatedAt,
	}
}

func MapAddressDetailToAddressRes(row db.GetAddressByIDAndCustomerRow) *AddressRes {
	return &AddressRes{
		AddressID:     row.AddressID,
		CustomerID:    row.CustomerID,
		Title:         utils.NullStringToStringPtr(row.Title),
		RecipientName: row.RecipientName,
		AddressLine:   row.AddressLine,
		ProvinceID:    utils.Int32ToInt32Ptr(row.ProvinceID),
		Subdistrict:   row.Subdistrict,
		District:      row.District,
		Province:      row.Province,
		ZipCode:       row.ZipCode,
		PhoneNumber:   row.PhoneNumber,
		IsDefault:     utils.ParseBool(row.IsDefault),
		CreatedAt:     row.CreatedAt,
	}
}

func MapDbAddressesToAddressResList(dbAddrs []db.ListAddressesByCustomerIDRow) []*AddressRes {
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
		PhoneNumber:   req.PhoneNumber,
		IsDefault:     false,
		UpdatedAt:     utils.NullTimeNow(),
		UpdatedBy:     utils.Int32PtrToNullInt32(updatedBy),
		AddressID:     addressID,
		CustomerID:    customerID,
	}
}
