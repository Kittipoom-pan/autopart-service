package entitie

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbCustomerToCustomerRes(dbCustomer db.GetCustomerRow) *CustomerRes {
	return &CustomerRes{
		ID:          uint32(dbCustomer.CustomerID),
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
		Username: dbCustomer.Username,
		Password: dbCustomer.Password.String,
	}
}

func MapCustomerToCustomerParam(customer *CustomerReq, createBy string) db.CreateCustomerParams {
	return db.CreateCustomerParams{
		FirstName:   utils.StringToNullString(customer.FirstName),
		LastName:    utils.StringToNullString(customer.LastName),
		Username:    customer.Username,
		Email:       customer.Email,
		Password:    utils.StringToNullString(customer.Password),
		BirthDate:   utils.NullTime(customer.BirthDate),
		PhoneNumber: utils.StringToNullString(customer.PhoneNumber),
		CreatedBy:   utils.StringToNullString(createBy),
		CreatedAt:   utils.NullTimeNow(),
	}
}

func MapDbCustomersToCustomerRes(dbCustomer db.ListCustomersRow) *CustomerRes {
	return &CustomerRes{
		ID:          uint32(dbCustomer.CustomerID),
		FirstName:   dbCustomer.FirstName.String,
		LastName:    dbCustomer.LastName.String,
		Username:    dbCustomer.Username,
		Email:       dbCustomer.Email,
		BirthDate:   dbCustomer.BirthDate.Time,
		PhoneNumber: dbCustomer.PhoneNumber.String,
	}
}

func MapUpdateCustomerParams(id int, customer *CustomerReq, updatedBy string) db.UpdateCustomerParams {
	return db.UpdateCustomerParams{
		CustomerID:  int32(id),
		FirstName:   utils.StringToNullString(customer.FirstName),
		LastName:    utils.StringToNullString(customer.LastName),
		Username:    customer.Username,
		Email:       customer.Email,
		BirthDate:   utils.NullTime(customer.BirthDate),
		PhoneNumber: utils.StringToNullString(customer.PhoneNumber),
		UpdatedBy:   utils.StringToNullString(updatedBy),
		UpdatedAt:   utils.NullTimeNow(),
	}
}

func MapUpdateCustomerIsActiveParams(id int, isActive bool, updatedBy string) db.UpdateCustomerIsActiveParams {
	return db.UpdateCustomerIsActiveParams{
		IsActive:   isActive,
		UpdatedBy:  utils.StringToNullString(updatedBy),
		UpdatedAt:  utils.NullTimeNow(),
		CustomerID: int32(id),
	}
}

func MapCustomerToLoginRes(customer *Customer, token string, expiresIn int32) *LoginResponse {
	userRes := &UserResponse{
		ID:       int32(customer.ID),
		Username: customer.Username,
		Role:     "",
	}

	return &LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		User:        *userRes,
	}
}
