package entitie

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbCustomerToCustomerRes(dbCustomer db.GetCustomerRow) *CustomerRes {
	return &CustomerRes{
		ID:          uint32(dbCustomer.CustomerID),
		Uuid:        string(dbCustomer.Uuid),
		FirstName:   dbCustomer.FirstName.String,
		LastName:    dbCustomer.LastName.String,
		Username:    dbCustomer.Username,
		Email:       dbCustomer.Email,
		BirthDate:   dbCustomer.BirthDate.Time,
		PhoneNumber: dbCustomer.PhoneNumber.String,
	}
}

func MapDbCustomerToCustomerEntity(dbCustomer db.GetCustomerByUsernameRow) *Customer {
	return &Customer{
		ID:       uint32(dbCustomer.CustomerID),
		Uuid:     string(dbCustomer.Uuid),
		Username: dbCustomer.Username,
		Password: dbCustomer.Password.String,
	}
}

func MapCustomerToCustomerParam(customer *CustomerReq, createBy *int) db.CreateCustomerParams {
	return db.CreateCustomerParams{
		Uuid:        utils.NewUUIDBytes(),
		FirstName:   utils.StringToNullString(customer.FirstName),
		LastName:    utils.StringToNullString(customer.LastName),
		Username:    customer.Username,
		Email:       customer.Email,
		Password:    utils.StringToNullString(customer.Password),
		BirthDate:   utils.NullTime(customer.BirthDate),
		PhoneNumber: utils.StringToNullString(customer.PhoneNumber),
		CreatedBy:   utils.IntNullToNullInt32(createBy),
		CreatedAt:   utils.NullTimeNow(),
	}
}

func MapDbCustomersToCustomerRes(dbCustomer db.ListCustomersRow) *CustomerRes {
	return &CustomerRes{
		ID:          uint32(dbCustomer.CustomerID),
		Uuid:        string(dbCustomer.Uuid),
		FirstName:   dbCustomer.FirstName.String,
		LastName:    dbCustomer.LastName.String,
		Username:    dbCustomer.Username,
		Email:       dbCustomer.Email,
		BirthDate:   dbCustomer.BirthDate.Time,
		PhoneNumber: dbCustomer.PhoneNumber.String,
	}
}

func MapUpdateCustomerParams(id int, customer *CustomerReq, updatedBy int) db.UpdateCustomerParams {
	return db.UpdateCustomerParams{
		CustomerID:  int32(id),
		FirstName:   utils.StringToNullString(customer.FirstName),
		LastName:    utils.StringToNullString(customer.LastName),
		Username:    customer.Username,
		Email:       customer.Email,
		BirthDate:   utils.NullTime(customer.BirthDate),
		PhoneNumber: utils.StringToNullString(customer.PhoneNumber),
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
