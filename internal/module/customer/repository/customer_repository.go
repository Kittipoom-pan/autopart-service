package repository

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
)

type CustomerRepositoryImpl struct {
	queries *db.Queries
}

type CustomerRepository interface {
	GetCustomerByID(id int) (*entitie.Customer, error)
	CreateCustomer(user *entitie.Customer) (int64, error)
}

func NewCustomerRepository(queries *db.Queries) CustomerRepository {
	return &CustomerRepositoryImpl{
		queries: queries,
	}
}

func (r *CustomerRepositoryImpl) GetCustomerByID(id int) (*entitie.Customer, error) {
	customer, err := r.queries.GetCustomer(context.Background(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	fmt.Sprintf(customer.Email)

	return &entitie.Customer{
		ID:          uint32(customer.CustomerID),
		FirstName:   customer.FirstName.String,
		LastName:    customer.LastName.String,
		Username:    customer.Username,
		Email:       customer.Email,
		Password:    customer.Password.String,
		BirthDate:   customer.BirthDate.Time,
		PhoneNumber: customer.PhoneNumber.String,
	}, nil
}

func (r *CustomerRepositoryImpl) CreateCustomer(user *entitie.Customer) (int64, error) {
	params := db.CreateCustomerParams{
		FirstName:   sql.NullString{String: user.FirstName, Valid: user.FirstName != ""},
		LastName:    sql.NullString{String: user.LastName, Valid: user.LastName != ""},
		Username:    user.Username,
		Email:       user.Email,
		Password:    sql.NullString{String: user.Password, Valid: user.Password != ""},
		BirthDate:   sql.NullTime{Time: user.BirthDate, Valid: !user.BirthDate.IsZero()},
		PhoneNumber: sql.NullString{String: user.PhoneNumber, Valid: user.PhoneNumber != ""},
		//CreatedAt:   user.CreatedAt,
		CreatedBy: sql.NullString{String: user.CreatedBy, Valid: user.CreatedBy != ""},
		UpdatedAt: sql.NullTime{Time: user.UpdatedAt, Valid: !user.UpdatedAt.IsZero()},
		UpdatedBy: sql.NullString{String: user.UpdatedBy, Valid: user.UpdatedBy != ""},
	}

	result, err := r.queries.CreateCustomer(context.Background(), params)
	if err != nil {
		return 0, err
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return customerID, nil
}
