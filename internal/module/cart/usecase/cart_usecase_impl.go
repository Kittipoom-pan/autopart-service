package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/usecase/validation"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
)

type cartUsecase struct {
	repo      repository.CartRepository
	txManager database.TxManager
	logger    zerolog.Logger
}

func NewCartUsecase(repo repository.CartRepository, txManager database.TxManager, logger zerolog.Logger) CartUsecase {
	return &cartUsecase{
		repo:      repo,
		txManager: txManager,
		logger:    logger,
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

	// ตรวจ part ก่อนเข้า tx จะได้ไม่สร้าง cart เปล่าจาก request ที่ผิด
	if err := u.repo.CheckPartExists(ctx, req.PartID); err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Part")
	}

	var (
		cartID      int32
		updatedCart *entity.CartRes
	)

	err := u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// หา active cart ของลูกค้า ถ้ายังไม่มี (ลูกค้าใหม่ หรือเพิ่ง checkout) ให้สร้างให้
		// statement นี้ล็อกแถว cart ไว้จนจบ tx ด้วย: AddItem ของ cart เดียวกันจึงเข้าคิวทีละ request
		id, err := u.repo.GetOrCreateActiveCartID(ctx, customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart")
		}
		cartID = id

		existingItem, err := u.repo.GetCartItemByCartAndPart(ctx, cartID, req.PartID)
		if err != nil && !dberr.IsRecordNotFound(err) {
			return helper.MapDBErrorToAPIError(err, "Cart item")
		}

		if existingItem != nil {
			newQuantity := existingItem.Quantity + req.Quantity
			if err := u.repo.UpdateCartItemQuantity(ctx, existingItem.CartItemID, cartID, newQuantity); err != nil {
				return helper.MapDBErrorToAPIError(err, "Cart item")
			}
		} else {
			if _, err := u.repo.CreateCartItem(ctx, cartID, req.PartID, req.Quantity); err != nil {
				return helper.MapDBErrorToAPIError(err, "Cart item")
			}
		}

		if err := u.repo.RecalcCartAmount(ctx, cartID); err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart")
		}

		updatedCart, err = u.repo.GetActiveCartByCustomerID(ctx, customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	u.logger.Info().Int32("cart_id", cartID).Int32("part_id", req.PartID).Msg("item added to cart successfully")
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
	err = u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := u.repo.UpdateCartItemQuantity(ctx, cartItemID, cart.CartID, req.Quantity); err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart item")
		}

		if err := u.repo.RecalcCartAmount(ctx, cart.CartID); err != nil {
			return helper.MapDBErrorToAPIError(err, "Cart")
		}

		updatedCart, err = u.repo.GetActiveCartByCustomerID(ctx, customerID)
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
