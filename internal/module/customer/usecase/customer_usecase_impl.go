package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
)

type CustomerUsecaseImpl struct {
	repo repository.CustomerRepository
}

func NewCustomerUsecase(repo repository.CustomerRepository) CustomerUsecase {
	return &CustomerUsecaseImpl{repo: repo}
}

func (u *CustomerUsecaseImpl) GetCustomerByID(ctx context.Context, id int) (*entitie.Customer, error) {
	customer, err := u.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (u *CustomerUsecaseImpl) CreateCustomer(ctx context.Context, user *entitie.CustomerReq) (int64, error) {
	customerID, err := u.repo.CreateCustomer(ctx, user)
	if err != nil {
		return 0, err
	}
	return customerID, nil
}
