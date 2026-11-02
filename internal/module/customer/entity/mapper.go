package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbCustomerToCustomerRes(dbCustomer db.GetCustomerRow) *CustomerRes {
	return &CustomerRes{
		ID:          int32(dbCustomer.CustomerID),
		FirstName:   dbCustomer.FirstName,
		LastName:    dbCustomer.LastName,
		Username:    dbCustomer.Username,
		Email:       dbCustomer.Email,
		BirthDate:   utils.NullTimeToDateStr(dbCustomer.BirthDate),
		PhoneNumber: dbCustomer.PhoneNumber,
	}
}

func MapDbCustomerToCustomerEntity(dbCustomer db.GetCustomerByUsernameRow) *Customer {
	return &Customer{
		ID:       int32(dbCustomer.CustomerID),
		Username: dbCustomer.Username,
		Password: dbCustomer.Password,
	}
}

func MapCustomerToCustomerParam(customer *CustomerReq, createBy *int) db.CreateCustomerParams {
	return db.CreateCustomerParams{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Username:    customer.Username,
		Email:       customer.Email,
		Password:    customer.Password,
		BirthDate:   utils.NullTime(customer.BirthDate),
		PhoneNumber: customer.PhoneNumber,
		CreatedBy:   utils.IntNullToNullInt32(createBy),
		CreatedAt:   utils.NowUTC(),
	}
}

func MapDbCustomersToCustomerRes(dbCustomer db.ListCustomersRow) *CustomerRes {
	return &CustomerRes{
		ID:          int32(dbCustomer.CustomerID),
		FirstName:   dbCustomer.FirstName,
		LastName:    dbCustomer.LastName,
		Username:    dbCustomer.Username,
		Email:       dbCustomer.Email,
		BirthDate:   utils.NullTimeToDateStr(dbCustomer.BirthDate),
		PhoneNumber: dbCustomer.PhoneNumber,
	}
}

func MapUpdateCustomerParams(id int, customer *CustomerReq, updatedBy int) db.UpdateCustomerParams {
	return db.UpdateCustomerParams{
		CustomerID:  int32(id),
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Username:    customer.Username,
		Email:       customer.Email,
		BirthDate:   utils.NullTime(customer.BirthDate),
		PhoneNumber: customer.PhoneNumber,
		UpdatedBy:   utils.IntToNullInt32(updatedBy),
		UpdatedAt:   utils.NullTimeNow(),
	}
}

func MapUpdateCustomerIsActiveParams(id int, isActive bool, updatedBy int) db.UpdateCustomerIsActiveParams {
	return db.UpdateCustomerIsActiveParams{
		IsActive:   isActive,
		UpdatedBy:  utils.IntToNullInt32(updatedBy),
		UpdatedAt:  utils.NullTimeNow(),
		CustomerID: int32(id),
	}
}

func MapCustomerToLoginRes(customer *Customer, token string, expiresIn int32) *LoginResponse {
	return &LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
	}
}
