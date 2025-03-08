package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
)

type CustomerUsecase interface {
	GetCustomerByID(ctx context.Context, id int) (*entitie.Customer, error)
	CreateCustomer(ctx context.Context, user *entitie.CustomerReq) (int64, error)
}
