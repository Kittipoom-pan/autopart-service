package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/usecase/validation"
	"github.com/rs/zerolog"
)

type addressUsecase struct {
	repo      repository.AddressRepository
	txManager database.TxManager
	logger    zerolog.Logger
}

func NewAddressUsecase(repo repository.AddressRepository, txManager database.TxManager, logger zerolog.Logger) AddressUsecase {
	return &addressUsecase{
		repo:      repo,
		txManager: txManager,
		logger:    logger,
	}
}

func (u *addressUsecase) GetCustomerAddresses(ctx context.Context, customerID int32) ([]*entity.AddressRes, error) {
	u.logger.Debug().Int32("customer_id", customerID).Msg("GetCustomerAddresses")

	addresses, err := u.repo.GetCustomerAddresses(ctx, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Address")
	}

	if addresses == nil {
		addresses = []*entity.AddressRes{}
	}

	return addresses, nil
}

func (u *addressUsecase) GetAddressByIDAndCustomer(ctx context.Context, addressID int32, customerID int32) (*entity.AddressRes, error) {
	u.logger.Debug().Int32("address_id", addressID).Int32("customer_id", customerID).Msg("GetAddressByIDAndCustomer")

	address, err := u.repo.GetAddressByIDAndCustomer(ctx, addressID, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Address")
	}

	return address, nil
}

func (u *addressUsecase) CreateAddress(ctx context.Context, customerID int32, req *entity.CreateAddressReq) (*entity.AddressRes, error) {
	u.logger.Debug().Int32("customer_id", customerID).Interface("req", req).Msg("CreateAddress")

	if err := validation.ValidateCreateAddressRequest(req); err != nil {
		u.logger.Warn().Err(err).Msg("create address validation failed")
		return nil, err
	}

	var createdAddress *entity.AddressRes
	err := u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if req.IsDefault {
			if err := u.repo.ClearDefaultAddresses(ctx, customerID, &customerID); err != nil {
				return helper.MapDBErrorToAPIError(err, "Address")
			}
		}

		id, err := u.repo.CreateAddress(ctx, req, customerID, &customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Address")
		}

		addr, err := u.repo.GetAddressByIDAndCustomer(ctx, int32(id), customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Address")
		}

		createdAddress = addr
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.Info().Int32("address_id", createdAddress.AddressID).Int32("customer_id", customerID).Msg("address created successfully")
	return createdAddress, nil
}

func (u *addressUsecase) UpdateAddress(ctx context.Context, addressID int32, customerID int32, req *entity.UpdateAddressReq) (*entity.AddressRes, error) {
	u.logger.Debug().Int32("address_id", addressID).Int32("customer_id", customerID).Interface("req", req).Msg("UpdateAddress")

	if err := validation.ValidateUpdateAddressRequest(addressID, req); err != nil {
		u.logger.Warn().Err(err).Msg("update address validation failed")
		return nil, err
	}

	// Verify address exists and belongs to customer
	if _, err := u.repo.GetAddressByIDAndCustomer(ctx, addressID, customerID); err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Address")
	}

	if err := u.repo.UpdateAddress(ctx, addressID, customerID, req, &customerID); err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Address")
	}

	updated, err := u.repo.GetAddressByIDAndCustomer(ctx, addressID, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Address")
	}

	u.logger.Info().Int32("address_id", addressID).Int32("customer_id", customerID).Msg("address updated successfully")
	return updated, nil
}

func (u *addressUsecase) DeleteAddress(ctx context.Context, addressID int32, customerID int32) error {
	u.logger.Debug().Int32("address_id", addressID).Int32("customer_id", customerID).Msg("DeleteAddress")

	// Verify address exists and belongs to customer
	if _, err := u.repo.GetAddressByIDAndCustomer(ctx, addressID, customerID); err != nil {
		return helper.MapDBErrorToAPIError(err, "Address")
	}

	if err := u.repo.DeleteAddress(ctx, addressID, customerID); err != nil {
		return helper.MapDBErrorToAPIError(err, "Address")
	}

	u.logger.Info().Int32("address_id", addressID).Int32("customer_id", customerID).Msg("address deleted successfully")
	return nil
}

func (u *addressUsecase) SetDefaultAddress(ctx context.Context, addressID int32, customerID int32) (*entity.AddressRes, error) {
	u.logger.Debug().Int32("address_id", addressID).Int32("customer_id", customerID).Msg("SetDefaultAddress")

	// Verify address exists and belongs to customer
	if _, err := u.repo.GetAddressByIDAndCustomer(ctx, addressID, customerID); err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Address")
	}

	var updatedAddress *entity.AddressRes
	err := u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := u.repo.ClearDefaultAddresses(ctx, customerID, &customerID); err != nil {
			return helper.MapDBErrorToAPIError(err, "Address")
		}

		if err := u.repo.SetDefaultAddress(ctx, addressID, customerID, &customerID); err != nil {
			return helper.MapDBErrorToAPIError(err, "Address")
		}

		addr, err := u.repo.GetAddressByIDAndCustomer(ctx, addressID, customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Address")
		}

		updatedAddress = addr
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.Info().Int32("address_id", addressID).Int32("customer_id", customerID).Msg("default address set successfully")
	return updatedAddress, nil
}
