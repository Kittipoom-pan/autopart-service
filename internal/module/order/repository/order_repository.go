package repository

import (
	"context"
	"time"

	"github.com/Kittipoom-pan/autopart-service/internal/module/order/entity"
)

type OrderRepository interface {
	GetOrderByID(ctx context.Context, id int) (*entity.OrderRes, error)
	GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entity.OrderRes, error)
	GetAllOrders(ctx context.Context) ([]*entity.OrderRes, error)
	CreateOrder(ctx context.Context, in *entity.CreateOrderInput, orderNumber string) (int64, error)
	NextOrderSequence(ctx context.Context, seqDate time.Time) (int64, error)
	CreateOrderItem(ctx context.Context, orderID int64, item *entity.CreateOrderItemInput) error
	CreateOrderShippingAddress(ctx context.Context, orderID int64, address *entity.CreateOrderShippingAddressInput) error
}
