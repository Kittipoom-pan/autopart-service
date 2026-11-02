package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/entity"
)

type CustomerPaymentMethodUsecase interface {
	GetCustomerPaymentMethods(ctx context.Context, customerID int32) ([]*entity.CustomerPaymentMethodRes, error)
	GetCustomerPaymentMethodByID(ctx context.Context, id int32, customerID int32) (*entity.CustomerPaymentMethodRes, error)
	CreateCustomerPaymentMethod(ctx context.Context, customerID int32, req *entity.CreateCustomerPaymentMethodReq) (*entity.CustomerPaymentMethodRes, error)
	UpdateCustomerPaymentMethod(ctx context.Context, id int32, customerID int32, req *entity.UpdateCustomerPaymentMethodReq) (*entity.CustomerPaymentMethodRes, error)
	DeleteCustomerPaymentMethod(ctx context.Context, id int32, customerID int32) error
	SetDefaultCustomerPaymentMethod(ctx context.Context, id int32, customerID int32) (*entity.CustomerPaymentMethodRes, error)
}
