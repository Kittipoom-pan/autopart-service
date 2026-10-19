package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/entity"
)

type CartUsecase interface {
	GetCart(ctx context.Context, customerID int32) (*entity.CartRes, error)
	AddItem(ctx context.Context, customerID int32, req *entity.AddCartItemReq) (*entity.CartRes, error)
	UpdateItemQuantity(ctx context.Context, customerID int32, cartItemID int32, req *entity.UpdateCartItemQuantityReq) (*entity.CartRes, error)
}
