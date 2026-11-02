package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/address/entity"
)

type AddressUsecase interface {
	GetCustomerAddresses(ctx context.Context, customerID int32) ([]*entity.AddressRes, error)
	GetAddressByIDAndCustomer(ctx context.Context, addressID int32, customerID int32) (*entity.AddressRes, error)
	CreateAddress(ctx context.Context, customerID int32, req *entity.CreateAddressReq) (*entity.AddressRes, error)
	UpdateAddress(ctx context.Context, addressID int32, customerID int32, req *entity.UpdateAddressReq) (*entity.AddressRes, error)
	DeleteAddress(ctx context.Context, addressID int32, customerID int32) error
	SetDefaultAddress(ctx context.Context, addressID int32, customerID int32) (*entity.AddressRes, error)
}
