package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	dberror "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entity"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type customerRepository struct {
	conn    *sql.DB
	queries *db.Queries
	logger  zerolog.Logger
}

func NewCustomerRepository(conn *sql.DB, logger zerolog.Logger) CustomerRepository {
	return &customerRepository{
		conn:    conn,
		queries: db.New(conn),
		logger:  logger,
	}
}

func (r *customerRepository) getQueries(ctx context.Context) *db.Queries {
	if tx, ok := database.GetTx(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *customerRepository) GetCustomerByID(ctx context.Context, id int) (*entity.CustomerRes, error) {
	customer, err := r.getQueries(ctx).GetCustomer(ctx, int32(id))
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
	result, err := r.getQueries(ctx).CreateCustomer(ctx, params)
	if err != nil {
		dbErr := dberror.HandleMySQLError(err)
		if dberror.IsDuplicateKey(dbErr) {
			r.logger.Warn().Err(dbErr).Str("customer", params.Username.String).Msg("duplicate key error")
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
	customers, err := r.getQueries(ctx).ListCustomers(ctx)
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
	result, err := r.getQueries(ctx).UpdateCustomer(ctx, params)
	if err != nil {
		dbErr := dberror.HandleMySQLError(err)
		if dberror.IsDuplicateKey(dbErr) {
			r.logger.Warn().Err(dbErr).Str("customer", params.Username.String).Msg("duplicate key error")
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

// DeleteCustomer ลบบัญชีลูกค้าถาวร (anonymize): ล็อกแถว -> ล้างข้อมูลส่วนตัว + set deleted_at
// -> ลบ address/token และล้างบัตรที่บันทึกไว้ ต้องเรียกภายใน transaction เท่านั้น
func (r *customerRepository) DeleteCustomer(ctx context.Context, customerID int, updatedBy int) error {
	q := r.getQueries(ctx)

	if _, err := q.LockCustomerForDelete(ctx, int32(customerID)); err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("customer_id", customerID).Msg("customer not found for delete")
			return dberror.NewDBError(dberror.ErrRecordNotFound, "customer not found")
		}
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to lock customer for delete")
		return dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	result, err := q.DeleteCustomer(ctx, db.DeleteCustomerParams{
		UpdatedByAdminID: utils.IntToNullInt32(updatedBy),
		CustomerID:       int32(customerID),
	})
	if err != nil {
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to anonymize customer")
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

	if err := q.DeleteCustomerAddresses(ctx, int32(customerID)); err != nil {
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to delete customer addresses")
		return dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	if err := q.DeleteCustomerPaymentMethods(ctx, int32(customerID)); err != nil {
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to anonymize customer payment methods")
		return dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	if err := q.DeleteCustomerTokens(ctx, int32(customerID)); err != nil {
		r.logger.Error().Err(err).Int("customer_id", customerID).Msg("failed to delete customer tokens")
		return dberror.NewDBError(dberror.ErrDatabaseInternal, err.Error())
	}

	return nil
}

func (r *customerRepository) GetCustomerByUsername(ctx context.Context, username string) (*entity.Customer, error) {
	customer, err := r.getQueries(ctx).GetCustomerByUsername(ctx, utils.StringToNullString(username))
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
