package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/entity"
)

type CustomerPaymentMethodRepository interface {
	GetCustomerPaymentMethods(ctx context.Context, customerID int32) ([]*entity.CustomerPaymentMethodRes, error)
	GetCustomerPaymentMethodByID(ctx context.Context, id int32, customerID int32) (*entity.CustomerPaymentMethodRes, error)
	CreateCustomerPaymentMethod(ctx context.Context, req *entity.CreateCustomerPaymentMethodReq, customerID int32, createdBy *int32) (int64, error)
	UpdateCustomerPaymentMethod(ctx context.Context, id int32, customerID int32, req *entity.UpdateCustomerPaymentMethodReq, updatedBy *int32) error
	DeleteCustomerPaymentMethod(ctx context.Context, id int32, customerID int32) error
	SetDefaultCustomerPaymentMethod(ctx context.Context, id int32, customerID int32, updatedBy *int32) error
	ClearDefaultCustomerPaymentMethods(ctx context.Context, customerID int32, updatedBy *int32) error
}
