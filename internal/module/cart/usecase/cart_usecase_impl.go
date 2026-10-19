package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/usecase/validation"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
)

type cartUsecase struct {
	repo   repository.CartRepository
	logger zerolog.Logger
}

func NewCartUsecase(repo repository.CartRepository, logger zerolog.Logger) CartUsecase {
	return &cartUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *cartUsecase) GetCart(ctx context.Context, customerID int32) (*entity.CartRes, error) {
	u.logger.Debug().Int32("customer_id", customerID).Msg("GetCart")

	cart, err := u.repo.GetActiveCartByCustomerID(ctx, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Cart")
	}

	return cart, nil
}

func (u *cartUsecase) AddItem(ctx context.Context, customerID int32, req *entity.AddCartItemReq) (*entity.CartRes, error) {
	u.logger.Debug().Int32("customer_id", customerID).Interface("req", req).Msg("AddItem")

	if err := validation.ValidateAddCartItemRequest(req); err != nil {
		u.logger.Warn().Err(err).Msg("add cart item validation failed")
		return nil, err
	}

	// Requires an existing cart
	cart, err := u.repo.GetActiveCartByCustomerID(ctx, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Cart")
	}

	// Verify part exists
	if err := u.repo.CheckPartExists(ctx, req.PartID); err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Part")
	}

	var updatedCart *entity.CartRes
	err = u.repo.WithTx(ctx, func(tx repository.CartRepository) error {
		existingItem, err := tx.GetCartItemByCartAndPart(ctx, cart.CartID, req.PartID)
		if err != nil && !dberr.IsRecordNotFound(err) {
			return helper.MapDBErrorToAPIError(err, "Cart item")
		}

		if existingItem != nil {
			newQuantity := existingItem.Quantity + req.Quantity
			if err := tx.UpdateCartItemQuantity(ctx, existingItem.CartItemID, cart.CartID, newQuantity); err != nil {
				return helper.MapDBErrorToAPIError(err, "Cart item")
			}
		} else {
			if _, err := tx.CreateCartItem(ctx, cart.CartID, req.PartID, req.Quantity); err != nil {
				return helper.MapDBErrorToAPIError(err, "Cart item")
			}
		}

		if err := tx.RecalcCartAmount(ctx, cart.CartID); err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart")
		}

		updatedCart, err = tx.GetActiveCartByCustomerID(ctx, customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.Info().Int32("cart_id", cart.CartID).Int32("part_id", req.PartID).Msg("item added to cart successfully")
	return updatedCart, nil
}

func (u *cartUsecase) UpdateItemQuantity(ctx context.Context, customerID int32, cartItemID int32, req *entity.UpdateCartItemQuantityReq) (*entity.CartRes, error) {
	u.logger.Debug().Int32("customer_id", customerID).Int32("cart_item_id", cartItemID).Interface("req", req).Msg("UpdateItemQuantity")

	if err := validation.ValidateUpdateCartItemQuantityRequest(cartItemID, req); err != nil {
		u.logger.Warn().Err(err).Msg("update cart item quantity validation failed")
		return nil, err
	}

	// Requires an existing cart
	cart, err := u.repo.GetActiveCartByCustomerID(ctx, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Cart")
	}

	// Verify cart item exists and belongs to this cart
	existingItem, err := u.repo.GetCartItemByID(ctx, cartItemID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Cart item")
	}

	if existingItem.CartID != cart.CartID {
		u.logger.Warn().Int32("item_cart_id", existingItem.CartID).Int32("active_cart_id", cart.CartID).Msg("cart item does not belong to active cart")
		return nil, customerror.NewNotFoundError("Cart item")
	}

	var updatedCart *entity.CartRes
	err = u.repo.WithTx(ctx, func(tx repository.CartRepository) error {
		if err := tx.UpdateCartItemQuantity(ctx, cartItemID, cart.CartID, req.Quantity); err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart item")
		}

		if err := tx.RecalcCartAmount(ctx, cart.CartID); err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart")
		}

		updatedCart, err = tx.GetActiveCartByCustomerID(ctx, customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.Info().Int32("cart_id", cart.CartID).Int32("cart_item_id", cartItemID).Int32("quantity", req.Quantity).Msg("cart item quantity updated successfully")
	return updatedCart, nil
}
