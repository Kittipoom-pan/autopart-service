package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/stock/entity"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type stockRepository struct {
	conn    *sql.DB
	queries *db.Queries
	logger  zerolog.Logger
}

func NewStockRepository(conn *sql.DB, logger zerolog.Logger) StockRepository {
	return &stockRepository{conn: conn, queries: db.New(conn), logger: logger}
}

func (r *stockRepository) getQueries(ctx context.Context) *db.Queries {
	if tx, ok := database.GetTx(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *stockRepository) DeductStock(ctx context.Context, partID int32, quantity int32, updatedBy *int32) error {
	result, err := r.getQueries(ctx).DeductPartStock(ctx, db.DeductPartStockParams{
		Quantity:         quantity,
		UpdatedByAdminID: utils.Int32PtrToNullInt32(updatedBy),
		PartID:           partID,
	})
	if err != nil {
		r.logger.Error().Err(err).Int32("part_id", partID).Msg("failed to deduct part stock")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	affected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int32("part_id", partID).Msg("failed to read rows affected")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}
	if affected == 0 {
		r.logger.Warn().
			Int32("part_id", partID).
			Int32("quantity", quantity).
			Msg("insufficient stock")
		return ErrInsufficientStock
	}

	return nil
}

func (r *stockRepository) CreateStockMovement(ctx context.Context, in *entity.StockMovementInput) error {
	var orderID sql.NullInt64
	if in.OrderID != nil {
		orderID = sql.NullInt64{Int64: *in.OrderID, Valid: true}
	}

	_, err := r.getQueries(ctx).CreateStockMovement(ctx, db.CreateStockMovementParams{
		PartID:           in.PartID,
		QuantityChange:   in.QuantityChange,
		Price:            in.Price,
		EventType:        in.EventType,
		Reason:           in.Reason,
		OrderID:          orderID,
		Remark:           utils.StringPtrToNullString(in.Remark),
		CreatedByAdminID: utils.Int32PtrToNullInt32(in.CreatedBy),
	})
	if err != nil {
		r.logger.Error().Err(err).Int32("part_id", in.PartID).Msg("failed to create stock movement")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}
	return nil
}
