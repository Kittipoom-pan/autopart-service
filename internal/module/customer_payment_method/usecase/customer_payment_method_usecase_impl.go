package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/usecase/validation"
	"github.com/rs/zerolog"
)

type customerPaymentMethodUsecase struct {
	repo      repository.CustomerPaymentMethodRepository
	txManager database.TxManager
	logger    zerolog.Logger
}

func NewCustomerPaymentMethodUsecase(repo repository.CustomerPaymentMethodRepository, txManager database.TxManager, logger zerolog.Logger) CustomerPaymentMethodUsecase {
	return &customerPaymentMethodUsecase{
		repo:      repo,
		txManager: txManager,
		logger:    logger,
	}
}

func (u *customerPaymentMethodUsecase) GetCustomerPaymentMethods(ctx context.Context, customerID int32) ([]*entity.CustomerPaymentMethodRes, error) {
	u.logger.Debug().Int32("customer_id", customerID).Msg("GetCustomerPaymentMethods")

	methods, err := u.repo.GetCustomerPaymentMethods(ctx, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
	}

	if methods == nil {
		methods = []*entity.CustomerPaymentMethodRes{}
	}

	return methods, nil
}

func (u *customerPaymentMethodUsecase) GetCustomerPaymentMethodByID(ctx context.Context, id int32, customerID int32) (*entity.CustomerPaymentMethodRes, error) {
	u.logger.Debug().Int32("id", id).Int32("customer_id", customerID).Msg("GetCustomerPaymentMethodByID")

	method, err := u.repo.GetCustomerPaymentMethodByID(ctx, id, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
	}

	return method, nil
}

func (u *customerPaymentMethodUsecase) CreateCustomerPaymentMethod(ctx context.Context, customerID int32, req *entity.CreateCustomerPaymentMethodReq) (*entity.CustomerPaymentMethodRes, error) {
	u.logger.Debug().Int32("customer_id", customerID).Interface("req", req).Msg("CreateCustomerPaymentMethod")

	if err := validation.ValidateCreateCustomerPaymentMethodRequest(req); err != nil {
		u.logger.Warn().Err(err).Msg("create customer payment method validation failed")
		return nil, err
	}

	var createdMethod *entity.CustomerPaymentMethodRes
	err := u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if req.IsDefault {
			if err := u.repo.ClearDefaultCustomerPaymentMethods(ctx, customerID, &customerID); err != nil {
				return helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
			}
		}

		id, err := u.repo.CreateCustomerPaymentMethod(ctx, req, customerID, &customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
		}

		m, err := u.repo.GetCustomerPaymentMethodByID(ctx, int32(id), customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
		}

		createdMethod = m
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.Info().Int32("customer_payment_method_id", createdMethod.CustomerPaymentMethodID).Int32("customer_id", customerID).Msg("customer payment method created successfully")
	return createdMethod, nil
}

func (u *customerPaymentMethodUsecase) UpdateCustomerPaymentMethod(ctx context.Context, id int32, customerID int32, req *entity.UpdateCustomerPaymentMethodReq) (*entity.CustomerPaymentMethodRes, error) {
	u.logger.Debug().Int32("id", id).Int32("customer_id", customerID).Interface("req", req).Msg("UpdateCustomerPaymentMethod")

	if err := validation.ValidateUpdateCustomerPaymentMethodRequest(id, req); err != nil {
		u.logger.Warn().Err(err).Msg("update customer payment method validation failed")
		return nil, err
	}

	// Verify payment method exists and belongs to customer
	if _, err := u.repo.GetCustomerPaymentMethodByID(ctx, id, customerID); err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
	}

	if err := u.repo.UpdateCustomerPaymentMethod(ctx, id, customerID, req, &customerID); err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
	}

	updated, err := u.repo.GetCustomerPaymentMethodByID(ctx, id, customerID)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
	}

	u.logger.Info().Int32("customer_payment_method_id", id).Int32("customer_id", customerID).Msg("customer payment method updated successfully")
	return updated, nil
}

func (u *customerPaymentMethodUsecase) DeleteCustomerPaymentMethod(ctx context.Context, id int32, customerID int32) error {
	u.logger.Debug().Int32("id", id).Int32("customer_id", customerID).Msg("DeleteCustomerPaymentMethod")

	// Verify payment method exists and belongs to customer
	if _, err := u.repo.GetCustomerPaymentMethodByID(ctx, id, customerID); err != nil {
		return helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
	}

	if err := u.repo.DeleteCustomerPaymentMethod(ctx, id, customerID); err != nil {
		return helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
	}

	u.logger.Info().Int32("customer_payment_method_id", id).Int32("customer_id", customerID).Msg("customer payment method deleted successfully")
	return nil
}

func (u *customerPaymentMethodUsecase) SetDefaultCustomerPaymentMethod(ctx context.Context, id int32, customerID int32) (*entity.CustomerPaymentMethodRes, error) {
	u.logger.Debug().Int32("id", id).Int32("customer_id", customerID).Msg("SetDefaultCustomerPaymentMethod")

	// Verify payment method exists and belongs to customer
	if _, err := u.repo.GetCustomerPaymentMethodByID(ctx, id, customerID); err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
	}

	var updatedMethod *entity.CustomerPaymentMethodRes
	err := u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := u.repo.ClearDefaultCustomerPaymentMethods(ctx, customerID, &customerID); err != nil {
			return helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
		}

		if err := u.repo.SetDefaultCustomerPaymentMethod(ctx, id, customerID, &customerID); err != nil {
			return helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
		}

		m, err := u.repo.GetCustomerPaymentMethodByID(ctx, id, customerID)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "CustomerPaymentMethod")
		}

		updatedMethod = m
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.Info().Int32("customer_payment_method_id", id).Int32("customer_id", customerID).Msg("default customer payment method set successfully")
	return updatedMethod, nil
}
