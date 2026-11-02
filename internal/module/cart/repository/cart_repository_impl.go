package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
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

func (r *cartRepository) getQueries(ctx context.Context) *db.Queries {
	if tx, ok := database.GetTx(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *cartRepository) GetActiveCartByCustomerID(ctx context.Context, customerID int32) (*entity.CartRes, error) {
	cart, err := r.getQueries(ctx).GetCartByCustomerID(ctx, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int32("customer_id", customerID).Msg("active cart not found in database")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "cart not found")
		}
		r.logger.Error().Err(err).Int32("customer_id", customerID).Msg("failed to get active cart from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	items, err := r.getQueries(ctx).ListCartItemsWithPart(ctx, cart.CartID)
	if err != nil {
		r.logger.Error().Err(err).Int32("cart_id", cart.CartID).Msg("failed to list cart items from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	itemEntities := entity.MapDbCartItemsWithPartToCartItemResList(items)
	return entity.MapDbCartToCartRes(cart, itemEntities), nil
}

func (r *cartRepository) GetCartItemByCartAndPart(ctx context.Context, cartID int32, partID int32) (*entity.CartItem, error) {
	item, err := r.getQueries(ctx).GetCartItemByCartAndPart(ctx, db.GetCartItemByCartAndPartParams{
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
	item, err := r.getQueries(ctx).GetCartItemByID(ctx, cartItemID)
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
	result, err := r.getQueries(ctx).CreateCartItem(ctx, db.CreateCartItemParams{
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
	_, err := r.getQueries(ctx).UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
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
	_, err := r.getQueries(ctx).RecalcCartSubtotal(ctx, cartID)
	if err != nil {
		r.logger.Error().Err(err).Int32("cart_id", cartID).Msg("failed to recalc cart amount")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return nil
}

var ErrCartAlreadyCheckedOut = errors.New("cart already checked out")

func (r *cartRepository) MarkCartCheckedOut(ctx context.Context, cartID int32) error {
	res, err := r.getQueries(ctx).UpdateCartCheckoutByID(ctx, cartID)
	if err != nil {
		r.logger.Error().Err(err).Int32("cart_id", cartID).Msg("failed to mark cart as checked out")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	affected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int32("cart_id", cartID).Msg("failed to read rows affected")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}
	if affected == 0 {
		r.logger.Warn().Int32("cart_id", cartID).Msg("cart already checked out")
		return ErrCartAlreadyCheckedOut
	}

	return nil
}

func (r *cartRepository) CheckPartExists(ctx context.Context, partID int32) error {
	_, err := r.getQueries(ctx).GetPartByID(ctx, partID)
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
