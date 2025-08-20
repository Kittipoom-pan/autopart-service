package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
)

type CustomerRepository interface {
	GetCustomerByID(ctx context.Context, id int) (*entitie.CustomerRes, error)
	GetCustomerByUsername(ctx context.Context, username string) (*entitie.Customer, error)
	CreateCustomer(ctx context.Context, customer *entitie.CustomerReq) (int64, error)
	GetAllCustomers(ctx context.Context) ([]*entitie.CustomerRes, error)
	UpdateCustomer(ctx context.Context, id int, customer *entitie.CustomerReq) error
	DeleteCustomer(ctx context.Context, id int) error
}
