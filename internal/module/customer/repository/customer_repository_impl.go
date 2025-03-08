package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
)

type CustomerRepo struct {
	queries *db.Queries
}

func NewCustomerRepository(queries *db.Queries) CustomerRepository {
	return &CustomerRepo{
		queries: queries,
	}
}

func (r *CustomerRepo) GetCustomerByID(ctx context.Context, id int) (*entitie.Customer, error) {
	customer, err := r.queries.GetCustomer(ctx, int32(id))
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

func mapCustomerToParams(customer *entitie.CustomerReq) db.CreateCustomerParams {
	return db.CreateCustomerParams{
		FirstName:   sql.NullString{String: customer.FirstName, Valid: customer.FirstName != ""},
		LastName:    sql.NullString{String: customer.LastName, Valid: customer.LastName != ""},
		Username:    customer.Username,
		Email:       customer.Email,
		Password:    sql.NullString{String: customer.Password, Valid: customer.Password != ""},
		BirthDate:   sql.NullTime{Time: customer.BirthDate, Valid: !customer.BirthDate.IsZero()},
		PhoneNumber: sql.NullString{String: customer.PhoneNumber, Valid: customer.PhoneNumber != ""},
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

func (r *CustomerRepo) CreateCustomer(ctx context.Context, customer *entitie.CustomerReq) (int64, error) {
	params := mapCustomerToParams(customer)

	result, err := r.queries.CreateCustomer(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to create customer: %w", err)
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("repository: failed to retrieve customer ID: %w", err)
	}

	return customerID, nil
}
