package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbCustomerPaymentMethodToRes(dbMethod db.CustomerPaymentMethod) *CustomerPaymentMethodRes {
	return &CustomerPaymentMethodRes{
		CustomerPaymentMethodID: dbMethod.CustomerPaymentMethodID,
		CustomerID:              dbMethod.CustomerID,
		AccountName:             utils.NullStringToStringPtr(dbMethod.AccountName),
		PaymentMethodTypeID:     dbMethod.PaymentMethodTypeID,
		CardToken:               utils.NullStringToStringPtr(dbMethod.CardToken),
		IsDefault:               utils.ParseBool(dbMethod.IsDefault),
		CreatedAt:               dbMethod.CreatedAt,
		UpdatedAt:               utils.NullTimeToTimePtr(dbMethod.UpdatedAt),
	}
}

func MapDbCustomerPaymentMethodsToResList(dbMethods []db.CustomerPaymentMethod) []*CustomerPaymentMethodRes {
	items := make([]*CustomerPaymentMethodRes, 0, len(dbMethods))
	for _, m := range dbMethods {
		items = append(items, MapDbCustomerPaymentMethodToRes(m))
	}
	return items
}

func MapCreateCustomerPaymentMethodReqToParams(req *CreateCustomerPaymentMethodReq, customerID int32, createdBy *int32) db.CreateCustomerPaymentMethodParams {
	return db.CreateCustomerPaymentMethodParams{
		CustomerID:          customerID,
		AccountName:         utils.StringPtrToNullString(req.AccountName),
		PaymentMethodTypeID: req.PaymentMethodTypeID,
		CardToken:           utils.StringPtrToNullString(req.CardToken),
		IsDefault:           req.IsDefault,
		CreatedByAdminID:    utils.Int32PtrToNullInt32(createdBy),
	}
}

func MapUpdateCustomerPaymentMethodReqToParams(req *UpdateCustomerPaymentMethodReq, id int32, customerID int32, updatedBy *int32) db.UpdateCustomerPaymentMethodParams {
	return db.UpdateCustomerPaymentMethodParams{
		AccountName:             utils.StringPtrToNullString(req.AccountName),
		PaymentMethodTypeID:     req.PaymentMethodTypeID,
		CardToken:               utils.StringPtrToNullString(req.CardToken),
		UpdatedByAdminID:        utils.Int32PtrToNullInt32(updatedBy),
		CustomerPaymentMethodID: id,
		CustomerID:              customerID,
	}
}
