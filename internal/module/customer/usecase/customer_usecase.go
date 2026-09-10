package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entity"
)

type CustomerUsecase interface {
	GetCustomerByID(ctx context.Context, id int) (*entity.CustomerRes, error)
	GetAllCustomers(ctx context.Context) ([]*entity.CustomerRes, error)
	CreateCustomer(ctx context.Context, user *entity.CustomerReq) (int64, error)
	UpdateCustomer(ctx context.Context, customerID int, user *entity.CustomerReq, userID int) error
	DeleteCustomer(ctx context.Context, id int, userId int) error
}
