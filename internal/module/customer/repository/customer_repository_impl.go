package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type customerRepository struct {
	queries *db.Queries
	logger  zerolog.Logger
}

func NewCustomerRepository(queries *db.Queries) CustomerRepository {
	return &customerRepository{
		queries: queries,
		logger:  log.With().Str("component", "customer_repository").Logger(),
	}
}

func (r *customerRepository) GetCustomerByID(ctx context.Context, id int) (*entitie.Customer, error) {
	customer, err := r.queries.GetCustomer(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("customer_id", id).Msg("customer not found in database")
			return nil, customerror.NewNotFoundError("Customer")
		}
		r.logger.Error().Err(err).Int("customer_id", id).Msg("failed to get customer from database")
		return nil, customerror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	return entitie.MapDbCustomerToCustomerEntity(customer), nil
}

func (r *customerRepository) CreateCustomer(ctx context.Context, customer *entitie.CustomerReq) (int64, error) {
	params := entitie.MapCustomerToCustomerParam(customer, "SYSTEM")

	result, err := r.queries.CreateCustomer(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			r.logger.Warn().Err(mysqlErr).Str("customer", customer.Username).Msg("duplicate key error")
			return 0, customerror.NewAPIError(common.StatusConflict, "Username or email or phone number already exists")
		}
		r.logger.Error().Err(err).Msg("failed to create customer in database")
		return 0, customerror.NewAPIError(common.StatusError, "repository: failed to create customer")
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve customer ID from database")
		return 0, customerror.NewAPIError(common.StatusError, "repository: failed to retrieve customer ID")
	}

	return customerID, nil
}

func (r *customerRepository) GetAllCustomers(ctx context.Context) ([]*entitie.Customer, error) {
	customers, err := r.queries.ListCustomers(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list customers from database")
		return nil, customerror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	var customerEntities []*entitie.Customer
	for _, customer := range customers {
		customerEntities = append(customerEntities, entitie.MapDbCustomersToCustomerEntity(customer))
	}
	return customerEntities, nil
}

func (r *customerRepository) UpdateCustomer(ctx context.Context, id int, customer *entitie.CustomerReq) error {
	params := entitie.MapUpdateCustomerParams(id, customer, "SYSTEM")

	result, err := r.queries.UpdateCustomer(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			r.logger.Warn().Err(mysqlErr).Str("customer", customer.Username).Msg("duplicate key error")
			return customerror.NewAPIError(common.StatusConflict, "Username or email or phone number already exists")
		}
		r.logger.Error().Err(err).Int("customer_id", id).Msg("failed to update customer in database")
		return customerror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("customer_id", id).Msg("failed to get rows affected")
		return customerror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("customer_id", id).Msg("customer not found for update")
		return customerror.NewNotFoundError("Customer")
	}

	return nil
}

func (r *customerRepository) DeleteCustomer(ctx context.Context, id int) error {
	params := entitie.MapUpdateCustomerIsActiveParams(id, false, "SYSTEM")

	result, err := r.queries.UpdateCustomerIsActive(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int("customer_id", id).Msg("failed to delete customer in database")
		return customerror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("customer_id", id).Msg("failed to get rows affected")
		return customerror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("customer_id", id).Msg("customer not found for delete")
		return customerror.NewNotFoundError("Customer")
	}
	return nil
}
