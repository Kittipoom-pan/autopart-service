package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CustomerUse struct {
	repo   repository.CustomerRepository
	logger zerolog.Logger
}

func NewCustomerUsecase(repo repository.CustomerRepository) CustomerUsecase {
	return &CustomerUse{
		repo:   repo,
		logger: log.With().Str("component", "customer_usecase").Logger(),
	}
}

func (u *CustomerUse) GetCustomerByID(ctx context.Context, id int) (*entitie.Customer, error) {
	customer, err := u.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, customerror.NewAPIError(common.StatusError, "Failed to get customer: "+err.Error())
	}
	if customer == nil {
		return nil, customerror.NewNotFoundError("Customer")
	}

	return customer, nil
}

func (u *CustomerUse) CreateCustomer(ctx context.Context, user *entitie.CustomerReq) (int64, error) {
	customerID, err := u.repo.CreateCustomer(ctx, user)
	if err != nil {
		return 0, customerror.NewAPIError(common.StatusError, "Failed to create customer: "+err.Error())
	}
	return customerID, nil
}
