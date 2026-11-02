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

func MapDbCustomerPaymentMethodToRes(dbMethod db.CustomerPaymentMethod) *CustomerPaymentMethodRes {
	return &CustomerPaymentMethodRes{
		CustomerPaymentMethodID: dbMethod.CustomerPaymentMethodID,
		CustomerID:              dbMethod.CustomerID,
		AccountName:             utils.NullStringToStringPtr(dbMethod.AccountName),
		PaymentMethodTypeID:     dbMethod.PaymentMethodTypeID,
		CardToken:               utils.NullStringToStringPtr(dbMethod.CardToken),
		IsDefault:               parseBool(dbMethod.IsDefault),
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
		CreatedBy:           utils.Int32PtrToNullInt32(createdBy),
	}
}

func MapUpdateCustomerPaymentMethodReqToParams(req *UpdateCustomerPaymentMethodReq, id int32, customerID int32, updatedBy *int32) db.UpdateCustomerPaymentMethodParams {
	return db.UpdateCustomerPaymentMethodParams{
		AccountName:             utils.StringPtrToNullString(req.AccountName),
		PaymentMethodTypeID:     req.PaymentMethodTypeID,
		CardToken:               utils.StringPtrToNullString(req.CardToken),
		UpdatedBy:               utils.Int32PtrToNullInt32(updatedBy),
		CustomerPaymentMethodID: id,
		CustomerID:              customerID,
	}
}
