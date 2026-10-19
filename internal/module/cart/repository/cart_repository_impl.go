package repository

import (
	"context"
	"database/sql"

	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/entity"
	"github.com/rs/zerolog"
)

type cartRepository struct {
	conn    *sql.DB
	queries *db.Queries
	logger  zerolog.Logger
}

func NewCartRepository(conn *sql.DB, logger zerolog.Logger) CartRepository {
	return &cartRepository{conn: conn, queries: db.New(conn), logger: logger}
}

func (r *cartRepository) GetActiveCartByCustomerID(ctx context.Context, customerID int32) (*entity.CartRes, error) {
	cart, err := r.queries.GetCartByCustomerID(ctx, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int32("customer_id", customerID).Msg("active cart not found in database")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "cart not found")
		}
		r.logger.Error().Err(err).Int32("customer_id", customerID).Msg("failed to get active cart from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	items, err := r.queries.ListCartItemsWithPart(ctx, cart.CartID)
	if err != nil {
		r.logger.Error().Err(err).Int32("cart_id", cart.CartID).Msg("failed to list cart items from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	itemEntities := entity.MapDbCartItemsWithPartToCartItemResList(items)
	return entity.MapDbCartToCartRes(cart, itemEntities), nil
}

func (r *cartRepository) GetCartItemByCartAndPart(ctx context.Context, cartID int32, partID int32) (*entity.CartItem, error) {
	item, err := r.queries.GetCartItemByCartAndPart(ctx, db.GetCartItemByCartAndPartParams{
		CartID: cartID,
		PartID: partID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int32("cart_id", cartID).Int32("part_id", partID).Msg("cart item not found")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "cart item not found")
		}
		r.logger.Error().Err(err).Int32("cart_id", cartID).Int32("part_id", partID).Msg("failed to get cart item by cart and part")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbCartItemToCartItemEntity(item), nil
}

func (r *cartRepository) GetCartItemByID(ctx context.Context, cartItemID int32) (*entity.CartItem, error) {
	item, err := r.queries.GetCartItemByID(ctx, cartItemID)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int32("cart_item_id", cartItemID).Msg("cart item not found")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "cart item not found")
		}
		r.logger.Error().Err(err).Int32("cart_item_id", cartItemID).Msg("failed to get cart item by id")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbCartItemToCartItemEntity(item), nil
}

func (r *cartRepository) CreateCartItem(ctx context.Context, cartID int32, partID int32, quantity int32) (int64, error) {
	result, err := r.queries.CreateCartItem(ctx, db.CreateCartItemParams{
		CartID:   cartID,
		PartID:   partID,
		Quantity: quantity,
	})
	if err != nil {
		dbErr := dberr.HandleMySQLError(err)
		r.logger.Error().Err(err).Int32("cart_id", cartID).Int32("part_id", partID).Msg("failed to create cart item in database")
		return 0, dbErr
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve cart item id from database")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return id, nil
}

func (r *cartRepository) UpdateCartItemQuantity(ctx context.Context, cartItemID int32, cartID int32, quantity int32) error {
	_, err := r.queries.UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
		Quantity:   quantity,
		CartItemID: cartItemID,
		CartID:     cartID,
	})
	if err != nil {
		r.logger.Error().Err(err).Int32("cart_item_id", cartItemID).Int32("cart_id", cartID).Msg("failed to update cart item quantity")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}

func (r *cartRepository) RecalcCartAmount(ctx context.Context, cartID int32) error {
	_, err := r.queries.RecalcCartSubtotal(ctx, cartID)
	if err != nil {
		r.logger.Error().Err(err).Int32("cart_id", cartID).Msg("failed to recalc cart amount")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}

func (r *cartRepository) CheckPartExists(ctx context.Context, partID int32) error {
	_, err := r.queries.GetPartByID(ctx, partID)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int32("part_id", partID).Msg("part not found")
			return dberr.NewDBError(dberr.ErrRecordNotFound, "part not found")
		}
		r.logger.Error().Err(err).Int32("part_id", partID).Msg("failed to get part by id")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}

func (r *cartRepository) WithTx(ctx context.Context, fn func(CartRepository) error) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return dberr.NewDBError(dberr.ErrDatabaseInternal, "begin tx: "+err.Error())
	}

	txRepo := &cartRepository{conn: r.conn, queries: r.queries.WithTx(tx), logger: r.logger}

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
