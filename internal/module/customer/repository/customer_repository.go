package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entity"
)

type CustomerRepository interface {
	GetCustomerByID(ctx context.Context, id int) (*entity.CustomerRes, error)
	GetCustomerByUsername(ctx context.Context, username string) (*entity.Customer, error)
	CreateCustomer(ctx context.Context, customer *entity.CustomerReq, createdBy *int) (int64, error)
	GetAllCustomers(ctx context.Context) ([]*entity.CustomerRes, error)
	UpdateCustomer(ctx context.Context, customerID int, customer *entity.CustomerReq, updatedBy int) error
	DeleteCustomer(ctx context.Context, customerID int, updatedBy int) error
}
