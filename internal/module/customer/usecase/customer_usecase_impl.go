package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
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

func (u *customerUsecase) GetCustomerByID(ctx context.Context, id int) (*entitie.CustomerRes, error) {
	u.logger.Info().Int("customer_id", id).Msg("GetCustomerByID started")

	customer, err := u.repo.GetCustomerByID(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Int("customer_id", id).Msg("Failed to get customer from repository")
		return nil, err
	}

	u.logger.Info().Int("customer_id", id).Msg("GetCustomerByID completed successfully")
	return customer, nil
}

func (u *customerUsecase) GetAllCustomers(ctx context.Context) ([]*entitie.CustomerRes, error) {
	customers, err := u.repo.GetAllCustomers(ctx)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to get all customers from repository")
		return nil, err
	}
	return customers, nil
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, customer *entitie.CustomerReq) (int64, error) {
	hashedPassword, err := auth.HashPassword(customer.Password)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to hash password")
		return 0, err
	}

	customer.Password = hashedPassword
	customerID, err := u.repo.CreateCustomer(ctx, customer)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to create customer")
		return 0, err
	}
	u.logger.Info().Int64("customer_id", customerID).Msg("Customer created successfully")
	return customerID, nil
}

func (u *customerUsecase) UpdateCustomer(ctx context.Context, id int, user *entitie.CustomerReq) error {
	u.logger.Info().Int("customer_id", id).Msg("UpdateCustomer started")

	err := u.repo.UpdateCustomer(ctx, id, user)
	if err != nil {
		u.logger.Error().Err(err).Int("customer_id", id).Msg("Failed to update customer in repository")
		return err
	}

	u.logger.Info().Int("customer_id", id).Msg("UpdateCustomer completed successfully")
	return nil
}

func (u *customerUsecase) DeleteCustomer(ctx context.Context, id int) error {
	u.logger.Info().Int("customer_id", id).Msg("DeleteCustomer started")

	err := u.repo.DeleteCustomer(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Int("customer_id", id).Msg("Failed to delete customer in repository")
		return err
	}

	u.logger.Info().Int("customer_id", id).Msg("DeleteCustomer completed successfully")
	return nil
}
