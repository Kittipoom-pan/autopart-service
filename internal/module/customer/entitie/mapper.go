package entitie

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbCustomerToCustomerEntity(dbCustomer db.GetCustomerRow) *Customer {
	return &Customer{
		ID:          uint32(dbCustomer.CustomerID),
		FirstName:   dbCustomer.FirstName.String,
		LastName:    dbCustomer.LastName.String,
		Username:    dbCustomer.Username,
		Email:       dbCustomer.Email,
		BirthDate:   dbCustomer.BirthDate.Time,
		PhoneNumber: dbCustomer.PhoneNumber.String,
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

func MapDbCustomersToCustomerEntity(dbCustomer db.ListCustomersRow) *Customer {
	return &Customer{
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
