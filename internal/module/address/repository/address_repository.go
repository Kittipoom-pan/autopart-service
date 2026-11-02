package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/address/entity"
)

type AddressRepository interface {
	GetCustomerAddresses(ctx context.Context, customerID int32) ([]*entity.AddressRes, error)
	GetAddressByIDAndCustomer(ctx context.Context, addressID int32, customerID int32) (*entity.AddressRes, error)
	GetSubdistrictDetailByID(ctx context.Context, subdistrictID int32) (*entity.SubdistrictDetail, error)
	CreateAddress(ctx context.Context, req *entity.CreateAddressReq, customerID int32, createdBy *int32) (int64, error)
	UpdateAddress(ctx context.Context, addressID int32, customerID int32, req *entity.UpdateAddressReq, updatedBy *int32) error
	DeleteAddress(ctx context.Context, addressID int32, customerID int32) error
	SetDefaultAddress(ctx context.Context, addressID int32, customerID int32, updatedBy *int32) error
	ClearDefaultAddresses(ctx context.Context, customerID int32, updatedBy *int32) error
}
