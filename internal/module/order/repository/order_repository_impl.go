package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/entity"
	"github.com/rs/zerolog"
)

type orderRepository struct {
	conn    *sql.DB
	queries *db.Queries
	logger  zerolog.Logger
}

func NewOrderRepository(conn *sql.DB, logger zerolog.Logger) OrderRepository {
	return &orderRepository{conn: conn, queries: db.New(conn), logger: logger}
}

func (r *orderRepository) getQueries(ctx context.Context) *db.Queries {
	if tx, ok := database.GetTx(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id int) (*entity.OrderRes, error) {
	order, err := r.getQueries(ctx).GetOrderByID(ctx, int64(id))
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
	order, err := r.getQueries(ctx).GetOrderByOrderNumber(ctx, orderNumber)
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
	orders, err := r.getQueries(ctx).ListOrders(ctx)
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

func (r *orderRepository) CreateOrder(ctx context.Context, order *entity.CreateOrderInput, orderNumber string) (int64, error) {
	r.logger.Info().Interface("order_req", order).Msg("creating order")
	params := entity.MapCreateOrderInputToParams(order, orderNumber)
	r.logger.Info().Interface("order_params", params).Msg("creating order")

	result, err := r.getQueries(ctx).CreateOrder(ctx, params)
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

func (r *orderRepository) NextOrderSequence(ctx context.Context, seqDate time.Time) (int64, error) {
	result, err := r.getQueries(ctx).NextOrderSequence(ctx, seqDate)
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

func (r *orderRepository) CreateOrderItem(ctx context.Context, orderID int64, item *entity.CreateOrderItemInput) error {
	params := entity.MapCreateOrderItemInputToParams(item)
	params.OrderID = orderID
	_, err := r.getQueries(ctx).CreateOrderItem(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to create order item")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}
	return nil
}

func (r *orderRepository) CreateOrderShippingAddress(ctx context.Context, orderID int64, address *entity.CreateOrderShippingAddressInput) error {
	params := entity.MapCreateOrderShippingAddressInputToParams(orderID, address)
	_, err := r.getQueries(ctx).CreateOrderShippingAddress(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to create order shipping address")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}
	return nil
}
