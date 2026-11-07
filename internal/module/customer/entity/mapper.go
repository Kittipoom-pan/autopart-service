package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbCustomerToCustomerRes(dbCustomer db.GetCustomerRow) *CustomerRes {
	return &CustomerRes{
		ID:          int32(dbCustomer.CustomerID),
		FirstName:   utils.NullStringToString(dbCustomer.FirstName),
		LastName:    utils.NullStringToString(dbCustomer.LastName),
		Username:    utils.NullStringToString(dbCustomer.Username),
		Email:       utils.NullStringToString(dbCustomer.Email),
		BirthDate:   utils.NullTimeToDateStr(dbCustomer.BirthDate),
		PhoneNumber: utils.NullStringToString(dbCustomer.PhoneNumber),
	}
}

func MapDbCustomerToCustomerEntity(dbCustomer db.GetCustomerByUsernameRow) *Customer {
	return &Customer{
		ID:       int32(dbCustomer.CustomerID),
		Username: utils.NullStringToString(dbCustomer.Username),
		Password: utils.NullStringToString(dbCustomer.Password),
	}
}

func MapCustomerToCustomerParam(customer *CustomerReq, createBy *int) db.CreateCustomerParams {
	return db.CreateCustomerParams{
		FirstName:        utils.StringToNullString(customer.FirstName),
		LastName:         utils.StringToNullString(customer.LastName),
		Username:         utils.StringToNullString(customer.Username),
		Email:            utils.StringToNullString(customer.Email),
		Password:         utils.StringToNullString(customer.Password),
		BirthDate:        utils.NullTime(customer.BirthDate),
		PhoneNumber:      utils.StringToNullString(customer.PhoneNumber),
		CreatedByAdminID: utils.IntNullToNullInt32(createBy),
		CreatedAt:        utils.NowUTC(),
	}
}

func MapDbCustomersToCustomerRes(dbCustomer db.ListCustomersRow) *CustomerRes {
	return &CustomerRes{
		ID:          int32(dbCustomer.CustomerID),
		FirstName:   utils.NullStringToString(dbCustomer.FirstName),
		LastName:    utils.NullStringToString(dbCustomer.LastName),
		Username:    utils.NullStringToString(dbCustomer.Username),
		Email:       utils.NullStringToString(dbCustomer.Email),
		BirthDate:   utils.NullTimeToDateStr(dbCustomer.BirthDate),
		PhoneNumber: utils.NullStringToString(dbCustomer.PhoneNumber),
	}
}

func MapUpdateCustomerParams(id int, customer *CustomerReq, updatedBy int) db.UpdateCustomerParams {
	return db.UpdateCustomerParams{
		CustomerID:       int32(id),
		FirstName:        utils.StringToNullString(customer.FirstName),
		LastName:         utils.StringToNullString(customer.LastName),
		Username:         utils.StringToNullString(customer.Username),
		Email:            utils.StringToNullString(customer.Email),
		BirthDate:        utils.NullTime(customer.BirthDate),
		PhoneNumber:      utils.StringToNullString(customer.PhoneNumber),
		UpdatedByAdminID: utils.IntToNullInt32(updatedBy),
		UpdatedAt:        utils.NullTimeNow(),
	}
}

func MapUpdateCustomerIsActiveParams(id int, isActive bool, updatedBy int) db.UpdateCustomerIsActiveParams {
	return db.UpdateCustomerIsActiveParams{
		IsActive:         isActive,
		UpdatedByAdminID: utils.IntToNullInt32(updatedBy),
		UpdatedAt:        utils.NullTimeNow(),
		CustomerID:       int32(id),
	}
}

func MapCustomerToLoginRes(customer *Customer, token string, expiresIn int32) *LoginResponse {
	return &LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}
}
