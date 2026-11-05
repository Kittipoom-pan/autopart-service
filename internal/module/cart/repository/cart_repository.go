package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/entity"
)

type CartRepository interface {
	GetActiveCartByCustomerID(ctx context.Context, customerID int32) (*entity.CartRes, error)
	GetCartItemByCartAndPart(ctx context.Context, cartID int32, partID int32) (*entity.CartItem, error)
	GetCartItemByID(ctx context.Context, cartItemID int32) (*entity.CartItem, error)
	CreateCartItem(ctx context.Context, cartID int32, partID int32, quantity int32) (int64, error)
	UpdateCartItemQuantity(ctx context.Context, cartItemID int32, cartID int32, quantity int32) error
	RecalcCartAmount(ctx context.Context, cartID int32) error
	MarkCartCheckedOut(ctx context.Context, cartID int32) error
	CheckPartExists(ctx context.Context, partID int32) error
	GetOrCreateActiveCartID(ctx context.Context, customerID int32) (int32, error)
}
