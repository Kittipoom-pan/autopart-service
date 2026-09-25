package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/usecase/validation"
	customererror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type customerUsecase struct {
	repo   repository.CustomerRepository
	logger zerolog.Logger
}

func NewCustomerUsecase(repo repository.CustomerRepository) CustomerUsecase {
	return &customerUsecase{
		repo:   repo,
		logger: log.With().Str("component", "customer_usecase").Logger(),
	}
}

func (u *customerUsecase) GetCustomerByID(ctx context.Context, id int) (*entity.CustomerRes, error) {
	u.logger.Debug().Int("customer_id", id).Msg("GetCustomerByID")
	customer, err := u.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Customer")
	}
	return customer, nil
}

func (u *customerUsecase) GetAllCustomers(ctx context.Context) ([]*entity.CustomerRes, error) {
	u.logger.Debug().Msg("GetAllCustomers")
	customers, err := u.repo.GetAllCustomers(ctx)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Customer")
	}
	return customers, nil
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, customer *entity.CustomerReq) (int64, error) {
	if err := validation.ValidateCustomerRequest(customer, false); err != nil {
		u.logger.Warn().Err(err).Msg("customer create validation failed")
		return 0, err
	}

	hashedPassword, err := auth.HashPassword(customer.Password)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to hash password")
		return 0, customererror.NewAPIError(common.StatusError, "failed to process password")
	}

	req := *customer
	req.Password = hashedPassword

	customerID, err := u.repo.CreateCustomer(ctx, &req, nil)
	if err != nil {
		return 0, helper.MapDBErrorToAPIError(err, "Customer")
	}

	u.logger.Info().Int64("customer_id", customerID).Msg("customer created successfully")
	return customerID, nil
}

func (u *customerUsecase) UpdateCustomer(ctx context.Context, customerID int, customer *entity.CustomerReq, userID int) error {
	if err := validation.ValidateCustomerRequest(customer, true); err != nil {
		u.logger.Warn().Err(err).Msg("customer update validation failed")
		return err
	}

	if err := u.repo.UpdateCustomer(ctx, customerID, customer, userID); err != nil {
		return helper.MapDBErrorToAPIError(err, "Customer")
	}

	u.logger.Info().Int("customer_id", customerID).Msg("customer updated successfully")
	return nil
}

func (u *customerUsecase) DeleteCustomer(ctx context.Context, id int, userId int) error {
	if err := u.repo.DeleteCustomer(ctx, id, userId); err != nil {
		return helper.MapDBErrorToAPIError(err, "Customer")
	}

	u.logger.Info().Int("customer_id", id).Msg("customer deleted successfully")
	return nil
}
