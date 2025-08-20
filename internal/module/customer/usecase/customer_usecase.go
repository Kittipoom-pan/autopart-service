package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
)

type CustomerUsecase interface {
	GetCustomerByID(ctx context.Context, id int) (*entitie.CustomerRes, error)
	GetAllCustomers(ctx context.Context) ([]*entitie.CustomerRes, error)
	CreateCustomer(ctx context.Context, user *entitie.CustomerReq) (int64, error)
	UpdateCustomer(ctx context.Context, id int, user *entitie.CustomerReq) error
	DeleteCustomer(ctx context.Context, id int) error
}
