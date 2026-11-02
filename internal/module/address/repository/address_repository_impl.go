package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/entity"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type addressRepository struct {
	conn    *sql.DB
	queries *db.Queries
	logger  zerolog.Logger
}

func NewAddressRepository(conn *sql.DB, logger zerolog.Logger) AddressRepository {
	return &addressRepository{
		conn:    conn,
		queries: db.New(conn),
		logger:  logger,
	}
}

func (r *addressRepository) getQueries(ctx context.Context) *db.Queries {
	if tx, ok := database.GetTx(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *addressRepository) GetCustomerAddresses(ctx context.Context, customerID int32) ([]*entity.AddressRes, error) {
	addresses, err := r.getQueries(ctx).ListAddressesByCustomerID(ctx, customerID)
	if err != nil {
		r.logger.Error().Err(err).Int32("customer_id", customerID).Msg("failed to get customer addresses")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbAddressesToAddressResList(addresses), nil
}

func (r *addressRepository) GetAddressByIDAndCustomer(ctx context.Context, addressID int32, customerID int32) (*entity.AddressRes, error) {
	address, err := r.getQueries(ctx).GetAddressByIDAndCustomer(ctx, db.GetAddressByIDAndCustomerParams{
		AddressID:  addressID,
		CustomerID: customerID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn().Int32("address_id", addressID).Int32("customer_id", customerID).Msg("address not found")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "address not found")
		}
		r.logger.Error().Err(err).Int32("address_id", addressID).Int32("customer_id", customerID).Msg("failed to get address by id")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapAddressDetailToAddressRes(address), nil
}

func (r *addressRepository) GetSubdistrictDetailByID(ctx context.Context, subdistrictID int32) (*entity.SubdistrictDetail, error) {
	detail, err := r.getQueries(ctx).GetSubdistrictDetailByID(ctx, subdistrictID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn().Int32("subdistrict_id", subdistrictID).Msg("subdistrict not found")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "subdistrict not found")
		}
		r.logger.Error().Err(err).Int32("subdistrict_id", subdistrictID).Msg("failed to get subdistrict detail")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return &entity.SubdistrictDetail{
		SubdistrictID: detail.SubdistrictID,
		Subdistrict:   detail.Subdistrict,
		District:      detail.District,
		ProvinceID:    detail.ProvinceID,
		Province:      detail.Province,
		ZipCode:       detail.ZipCode,
	}, nil
}

func (r *addressRepository) CreateAddress(ctx context.Context, req *entity.CreateAddressReq, customerID int32, createdBy *int32) (int64, error) {
	params := entity.MapCreateAddressReqToParams(req, customerID, createdBy)
	result, err := r.getQueries(ctx).CreateAddress(ctx, params)
	if err != nil {
		dbErr := dberr.HandleMySQLError(err)
		r.logger.Error().Err(err).Int32("customer_id", customerID).Msg("failed to create address")
		return 0, dbErr
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get last insert id for address")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return id, nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, addressID int32, customerID int32, req *entity.UpdateAddressReq, updatedBy *int32) error {
	params := entity.MapUpdateAddressReqToParams(req, addressID, customerID, updatedBy)
	if _, err := r.getQueries(ctx).UpdateAddress(ctx, params); err != nil {
		r.logger.Error().Err(err).Int32("address_id", addressID).Int32("customer_id", customerID).Msg("failed to update address")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}

func (r *addressRepository) DeleteAddress(ctx context.Context, addressID int32, customerID int32) error {
	result, err := r.getQueries(ctx).DeleteAddress(ctx, db.DeleteAddressParams{
		AddressID:  addressID,
		CustomerID: customerID,
	})
	if err != nil {
		r.logger.Error().Err(err).Int32("address_id", addressID).Int32("customer_id", customerID).Msg("failed to delete address")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get rows affected for delete address")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}
	if rowsAffected == 0 {
		return dberr.NewDBError(dberr.ErrRecordNotFound, "address not found")
	}

	return nil
}

func (r *addressRepository) SetDefaultAddress(ctx context.Context, addressID int32, customerID int32, updatedBy *int32) error {
	result, err := r.getQueries(ctx).SetDefaultAddress(ctx, db.SetDefaultAddressParams{
		UpdatedAt:  utils.NullTimeNow(),
		UpdatedBy:  utils.Int32PtrToNullInt32(updatedBy),
		AddressID:  addressID,
		CustomerID: customerID,
	})
	if err != nil {
		r.logger.Error().Err(err).Int32("address_id", addressID).Int32("customer_id", customerID).Msg("failed to set default address")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get rows affected for set default address")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}
	if rowsAffected == 0 {
		return dberr.NewDBError(dberr.ErrRecordNotFound, "address not found")
	}

	return nil
}

func (r *addressRepository) ClearDefaultAddresses(ctx context.Context, customerID int32, updatedBy *int32) error {
	if _, err := r.getQueries(ctx).ClearDefaultAddress(ctx, db.ClearDefaultAddressParams{
		UpdatedAt:  utils.NullTimeNow(),
		UpdatedBy:  utils.Int32PtrToNullInt32(updatedBy),
		CustomerID: customerID,
	}); err != nil {
		r.logger.Error().Err(err).Int32("customer_id", customerID).Msg("failed to clear default addresses")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}
