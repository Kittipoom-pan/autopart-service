package entitie

import (
	"database/sql"
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
)

func MapDbCustomerToCustomerEntity(dbCustomer db.GetCustomerRow) *Customer {
	return &Customer{
		ID:          uint32(dbCustomer.CustomerID),
		FirstName:   dbCustomer.FirstName.String,
		LastName:    dbCustomer.LastName.String,
		Username:    dbCustomer.Username,
		Email:       dbCustomer.Email,
		Password:    dbCustomer.Password.String,
		BirthDate:   dbCustomer.BirthDate.Time,
		PhoneNumber: dbCustomer.PhoneNumber.String,
	}
}

func MapCustomerToCustomerParam(customer *CustomerReq, createBy string) db.CreateCustomerParams {
	return db.CreateCustomerParams{
		FirstName:   sql.NullString{String: customer.FirstName, Valid: customer.FirstName != ""},
		LastName:    sql.NullString{String: customer.LastName, Valid: customer.LastName != ""},
		Username:    customer.Username,
		Email:       customer.Email,
		Password:    sql.NullString{String: customer.Password, Valid: customer.Password != ""},
		BirthDate:   sql.NullTime{Time: customer.BirthDate, Valid: !customer.BirthDate.IsZero()},
		PhoneNumber: sql.NullString{String: customer.PhoneNumber, Valid: customer.PhoneNumber != ""},
		CreatedBy:   sql.NullString{String: createBy, Valid: createBy != ""},
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

func MapDbCustomersToCustomerEntity(dbCustomer db.ListCustomersRow) *Customer {
	return &Customer{
		ID:          uint32(dbCustomer.CustomerID),
		FirstName:   dbCustomer.FirstName.String,
		LastName:    dbCustomer.LastName.String,
		Username:    dbCustomer.Username,
		Email:       dbCustomer.Email,
		Password:    dbCustomer.Password.String,
		BirthDate:   dbCustomer.BirthDate.Time,
		PhoneNumber: dbCustomer.PhoneNumber.String,
	}
}

func MapUpdateCustomerParams(id int, customer *CustomerReq, updatedBy string) db.UpdateCustomerParams {
	return db.UpdateCustomerParams{
		CustomerID:  int32(id),
		FirstName:   sql.NullString{String: customer.FirstName, Valid: customer.FirstName != ""},
		LastName:    sql.NullString{String: customer.LastName, Valid: customer.LastName != ""},
		Username:    customer.Username,
		Email:       customer.Email,
		BirthDate:   sql.NullTime{Time: customer.BirthDate, Valid: !customer.BirthDate.IsZero()},
		PhoneNumber: sql.NullString{String: customer.PhoneNumber, Valid: customer.PhoneNumber != ""},
		UpdatedBy:   sql.NullString{String: updatedBy, Valid: updatedBy != ""},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

func MapUpdateCustomerIsActiveParams(id int, isActive bool, updatedBy string) db.UpdateCustomerIsActiveParams {
	return db.UpdateCustomerIsActiveParams{
		IsActive: isActive,
		UpdatedBy: sql.NullString{
			String: updatedBy,
			Valid:  updatedBy != "",
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CustomerID: int32(id),
	}
}
