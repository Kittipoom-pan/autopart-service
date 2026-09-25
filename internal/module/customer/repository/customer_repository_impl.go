package repository

import (
	"context"
	"database/sql"

	dberror "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entity"
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

func (r *customerRepository) GetCustomerByID(ctx context.Context, id int) (*entity.CustomerRes, error) {
	customer, err := r.queries.GetCustomer(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("customer_id", id).Msg("customer not found in database")
			return nil, dberror.NewDBError(dberror.ErrRecordNotFound, "customer not found")
		}
		r.logger.Error().Err(err).Int("customer_id", id).Msg("failed to get customer from database")
		return nil, dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbCustomerToCustomerRes(customer), nil
}

func (r *customerRepository) CreateCustomer(ctx context.Context, customer *entity.CustomerReq, createdBy *int) (int64, error) {
	params := entity.MapCustomerToCustomerParam(customer, createdBy)
	result, err := r.queries.CreateCustomer(ctx, params)
	if err != nil {
		dbErr := dberror.HandleMySQLError(err)
		if dberror.IsDuplicateKey(dbErr) {
			r.logger.Warn().Err(dbErr).Str("customer", params.Username).Msg("duplicate key error")
			return 0, dbErr
		}
		r.logger.Error().Err(err).Msg("failed to create customer in database")
		return 0, dbErr
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve customer ID from database")
		return 0, dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	return customerID, nil
}

func (r *customerRepository) GetAllCustomers(ctx context.Context) ([]*entity.CustomerRes, error) {
	customers, err := r.queries.ListCustomers(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list customers from database")
		return nil, dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	customerEntities := make([]*entity.CustomerRes, 0, len(customers))
	for _, customer := range customers {
		customerEntities = append(customerEntities, entity.MapDbCustomersToCustomerRes(customer))
	}
	return customerEntities, nil
}

func (r *customerRepository) UpdateCustomer(ctx context.Context, customerID int, customer *entity.CustomerReq, updatedBy int) error {
	params := entity.MapUpdateCustomerParams(customerID, customer, updatedBy)
	result, err := r.queries.UpdateCustomer(ctx, params)
	if err != nil {
		dbErr := dberror.HandleMySQLError(err)
		if dberror.IsDuplicateKey(dbErr) {
			r.logger.Warn().Err(dbErr).Str("customer", params.Username).Msg("duplicate key error")
			return dbErr
		}
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to update customer in database")
		return dbErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to get rows affected")
		return dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("customer_id", customerID).Msg("customer not found for update")
		return dberror.NewDBError(dberror.ErrRecordNotFound, "customer not found")
	}

	return nil
}

func (r *customerRepository) DeleteCustomer(ctx context.Context, customerID int, updatedBy int) error {
	params := entity.MapUpdateCustomerIsActiveParams(customerID, false, updatedBy)
	result, err := r.queries.UpdateCustomerIsActive(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to delete customer in database")
		return dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to get rows affected")
		return dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("customer_id", customerID).Msg("customer not found for delete")
		return dberror.NewDBError(dberror.ErrRecordNotFound, "customer not found")
	}
	return nil
}

func (r *customerRepository) GetCustomerByUsername(ctx context.Context, username string) (*entity.Customer, error) {
	customer, err := r.queries.GetCustomerByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Str("username", username).Msg("customer not found in database")
			return nil, dberror.NewDBError(dberror.ErrRecordNotFound, "customer not found")
		}
		r.logger.Error().Err(err).Str("username", username).Msg("failed to get customer from database")
		return nil, dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbCustomerToCustomerEntity(customer), nil
}
