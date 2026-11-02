package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/entity"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type customerPaymentMethodRepository struct {
	conn    *sql.DB
	queries *db.Queries
	logger  zerolog.Logger
}

func NewCustomerPaymentMethodRepository(conn *sql.DB, logger zerolog.Logger) CustomerPaymentMethodRepository {
	return &customerPaymentMethodRepository{
		conn:    conn,
		queries: db.New(conn),
		logger:  logger,
	}
}

func (r *customerPaymentMethodRepository) getQueries(ctx context.Context) *db.Queries {
	if tx, ok := database.GetTx(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *customerPaymentMethodRepository) GetCustomerPaymentMethods(ctx context.Context, customerID int32) ([]*entity.CustomerPaymentMethodRes, error) {
	methods, err := r.getQueries(ctx).GetCustomerPaymentMethods(ctx, customerID)
	if err != nil {
		r.logger.Error().Err(err).Int32("customer_id", customerID).Msg("failed to get customer payment methods")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbCustomerPaymentMethodsToResList(methods), nil
}

func (r *customerPaymentMethodRepository) GetCustomerPaymentMethodByID(ctx context.Context, id int32, customerID int32) (*entity.CustomerPaymentMethodRes, error) {
	method, err := r.getQueries(ctx).GetCustomerPaymentMethodByID(ctx, db.GetCustomerPaymentMethodByIDParams{
		CustomerPaymentMethodID: id,
		CustomerID:              customerID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int32("customer_payment_method_id", id).Int32("customer_id", customerID).Msg("customer payment method not found")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "customer payment method not found")
		}
		r.logger.Error().Err(err).Int32("customer_payment_method_id", id).Int32("customer_id", customerID).Msg("failed to get customer payment method by id")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbCustomerPaymentMethodToRes(method), nil
}

func (r *customerPaymentMethodRepository) CreateCustomerPaymentMethod(ctx context.Context, req *entity.CreateCustomerPaymentMethodReq, customerID int32, createdBy *int32) (int64, error) {
	params := entity.MapCreateCustomerPaymentMethodReqToParams(req, customerID, createdBy)
	result, err := r.getQueries(ctx).CreateCustomerPaymentMethod(ctx, params)
	if err != nil {
		dbErr := dberr.HandleMySQLError(err)
		r.logger.Error().Err(err).Int32("customer_id", customerID).Msg("failed to create customer payment method")
		return 0, dbErr
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get last insert id for customer payment method")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return id, nil
}

func (r *customerPaymentMethodRepository) UpdateCustomerPaymentMethod(ctx context.Context, id int32, customerID int32, req *entity.UpdateCustomerPaymentMethodReq, updatedBy *int32) error {
	params := entity.MapUpdateCustomerPaymentMethodReqToParams(req, id, customerID, updatedBy)
	err := r.getQueries(ctx).UpdateCustomerPaymentMethod(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int32("customer_payment_method_id", id).Int32("customer_id", customerID).Msg("failed to update customer payment method")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}

func (r *customerPaymentMethodRepository) DeleteCustomerPaymentMethod(ctx context.Context, id int32, customerID int32) error {
	err := r.getQueries(ctx).DeleteCustomerPaymentMethod(ctx, db.DeleteCustomerPaymentMethodParams{
		CustomerPaymentMethodID: id,
		CustomerID:              customerID,
	})
	if err != nil {
		r.logger.Error().Err(err).Int32("customer_payment_method_id", id).Int32("customer_id", customerID).Msg("failed to delete customer payment method")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}

func (r *customerPaymentMethodRepository) SetDefaultCustomerPaymentMethod(ctx context.Context, id int32, customerID int32, updatedBy *int32) error {
	err := r.getQueries(ctx).SetDefaultCustomerPaymentMethod(ctx, db.SetDefaultCustomerPaymentMethodParams{
		CustomerPaymentMethodID: id,
		CustomerID:              customerID,
		UpdatedBy:               utils.Int32PtrToNullInt32(updatedBy),
	})
	if err != nil {
		r.logger.Error().Err(err).Int32("customer_payment_method_id", id).Int32("customer_id", customerID).Msg("failed to set default customer payment method")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}

func (r *customerPaymentMethodRepository) ClearDefaultCustomerPaymentMethods(ctx context.Context, customerID int32, updatedBy *int32) error {
	err := r.getQueries(ctx).ClearDefaultCustomerPaymentMethods(ctx, db.ClearDefaultCustomerPaymentMethodsParams{
		CustomerID: customerID,
		UpdatedBy:  utils.Int32PtrToNullInt32(updatedBy),
	})
	if err != nil {
		r.logger.Error().Err(err).Int32("customer_id", customerID).Msg("failed to clear default customer payment methods")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}
