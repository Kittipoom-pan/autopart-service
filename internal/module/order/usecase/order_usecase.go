package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/order/entity"
)

type OrderUsecase interface {
	GetOrderByID(ctx context.Context, id int) (*entity.OrderRes, error)
	GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entity.OrderRes, error)
	GetAllOrders(ctx context.Context) ([]*entity.OrderRes, error)
	CreateOrderByAdmin(ctx context.Context, req *entity.AdminCreateOrderReq, adminID int32) (string, error)
	CreateOrderByCustomer(ctx context.Context, req *entity.CustomerCreateOrderReq, customerID int32) (string, error)
}
