package repository

import (
	"context"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
)

type CustomerRepository interface {
	GetCustomerByID(ctx context.Context, id int) (*entitie.CustomerRes, error)
	GetCustomerByUsername(ctx context.Context, username string) (*entitie.Customer, error)
	CreateCustomer(ctx context.Context, params db.CreateCustomerParams) (int64, error)
	GetAllCustomers(ctx context.Context) ([]*entitie.CustomerRes, error)
	UpdateCustomer(ctx context.Context, params db.UpdateCustomerParams) error
	DeleteCustomer(ctx context.Context, params db.UpdateCustomerIsActiveParams) error
}
