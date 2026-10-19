package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/entity"
	"github.com/rs/zerolog"
)

type orderRepository struct {
	conn    *sql.DB // ใช้เปิด tx
	queries *db.Queries
	logger  zerolog.Logger
}

func NewOrderRepository(conn *sql.DB, logger zerolog.Logger) OrderRepository {
	return &orderRepository{conn: conn, queries: db.New(conn), logger: logger}
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id int) (*entity.OrderRes, error) {
	order, err := r.queries.GetOrderByID(ctx, int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("order_id", id).Msg("order not found in database")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "order not found")
		}
		r.logger.Error().Err(err).Int("order_id", id).Msg("failed to get order from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbOrderToOrderRes(order), nil
}

func (r *orderRepository) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entity.OrderRes, error) {
	order, err := r.queries.GetOrderByOrderNumber(ctx, orderNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Str("order_number", orderNumber).Msg("order not found in database")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "order not found")
		}
		r.logger.Error().Err(err).Str("order_number", orderNumber).Msg("failed to get order from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbOrderByOrderNumberToOrderRes(order), nil
}

func (r *orderRepository) GetAllOrders(ctx context.Context) ([]*entity.OrderRes, error) {
	orders, err := r.queries.ListOrders(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list orders from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	orderEntities := make([]*entity.OrderRes, 0, len(orders))
	for _, order := range orders {
		orderEntities = append(orderEntities, entity.MapDbOrdersToOrderRes(order))
	}
	return orderEntities, nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, in *entity.CreateOrderInput, orderNumber string) (int64, error) {
	params := entity.MapCreateOrderInputToParams(in, orderNumber)
	result, err := r.queries.CreateOrder(ctx, params)
	if err != nil {
		dbErr := dberr.HandleMySQLError(err)
		if dberr.IsDuplicateKey(dbErr) {
			r.logger.Warn().Err(dbErr).Str("order_number", params.OrderNumber).Msg("duplicate order number")
			return 0, dbErr
		}
		// if dberr.IsForeignKeyError(dbErr) {
		// 	r.logger.Warn().Err(dbErr).
		// 		Int32("customer_id", params.CustomerID).
		// 		Msg("foreign key constraint failed")
		// 	return 0, dbErr
		// }
		r.logger.Error().Err(err).Str("order_number", params.OrderNumber).Msg("failed to create order in database")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve order ID from database")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return orderID, nil
}

// fn ได้ repo ตัวใหม่ที่ผูกกับ tx; return error = rollback, nil = commit
func (r *orderRepository) WithTx(ctx context.Context, fn func(OrderRepository) error) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return dberr.NewDBError(dberr.ErrDatabaseInternal, "begin tx: "+err.Error())
	}

	txRepo := &orderRepository{conn: r.conn, queries: r.queries.WithTx(tx), logger: r.logger}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(txRepo); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.logger.Error().Err(rbErr).Msg("rollback failed")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return dberr.NewDBError(dberr.ErrDatabaseInternal, "commit tx: "+err.Error())
	}
	return nil
}

func (r *orderRepository) NextOrderSequence(ctx context.Context, seqDate time.Time) (int64, error) {
	result, err := r.queries.NextOrderSequence(ctx, seqDate)
	if err != nil {
		r.logger.Error().Err(err).Time("seq_date", seqDate).Msg("failed to get next order sequence")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, fmt.Sprintf("next order sequence: %v", err))
	}

	seq, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get last insert id")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, fmt.Sprintf("last insert id: %v", err))
	}

	return seq, nil
}
