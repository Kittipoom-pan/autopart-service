package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
)

type CustomerRepository interface {
	GetCustomerByID(ctx context.Context, id int) (*entitie.Customer, error)
	CreateCustomer(ctx context.Context, customer *entitie.CustomerReq) (int64, error)
}
