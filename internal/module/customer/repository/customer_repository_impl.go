package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CustomerRepo struct {
	queries *db.Queries
	logger  zerolog.Logger
}

func NewCustomerRepository(queries *db.Queries) CustomerRepository {
	return &CustomerRepo{
		queries: queries,
		logger:  log.With().Str("component", "customer_repository").Logger(),
	}
}

func (r *CustomerRepo) GetCustomerByID(ctx context.Context, id int) (*entitie.Customer, error) {
	customer, err := r.queries.GetCustomer(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("customer_id", id).Msg("customer not found in database")
			return nil, nil
		}
		r.logger.Error().Err(err).Int("customer_id", id).Msg("failed to get customer from database")
		return nil, customerror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}
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
		r.logger.Error().Err(err).Msg("failed to create customer in database")
		return 0, customerror.NewAPIError(common.StatusError, fmt.Sprintf("repository: failed to create customer: %v", err))
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve customer ID from database")
		return 0, customerror.NewAPIError(common.StatusError, fmt.Sprintf("repository: failed to retrieve customer ID: %v", err))
	}

	return customerID, nil
}
