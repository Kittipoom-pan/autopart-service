package usecase

import (
	"context"
	"errors"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type authUsecase struct {
	repo   repository.CustomerRepository
	logger zerolog.Logger
	cfg    *config.Config
}

func NewAuthUsecase(repo repository.CustomerRepository, cfg *config.Config) AuthUsecase {
	return &authUsecase{
		repo:   repo,
		cfg:    cfg,
		logger: log.With().Str("component", "customer_auth_usecase").Logger(),
	}
}

func (u *authUsecase) Login(ctx context.Context, request *entity.LoginRequest) (*entity.LoginResponse, error) {
	invalidCreds := customerror.NewAPIError(common.StatusUnauthorized, "Incorrect username or password")

	customer, err := u.repo.GetCustomerByUsername(ctx, request.Username)
	if err != nil {
		var notFound *customerror.NotFoundError
		if errors.As(err, &notFound) {
			u.logger.Warn().Str("username", request.Username).Msg("login failed: customer not found")
			return nil, invalidCreds
		}
		u.logger.Error().Err(err).Msg("failed to get customer for login")
		return nil, err
	}

	if err := auth.VerifyPassword(customer.Password, request.Password); err != nil {
		u.logger.Warn().Str("username", request.Username).Msg("login failed: invalid password")
		return nil, invalidCreds
	}

	token, err := auth.GenerateToken(customer.ID, auth.RoleCustomer, u.cfg)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to generate JWT token")
		return nil, customerror.NewAPIError(common.StatusError, "failed to generate token")
	}

	return entity.MapCustomerToLoginRes(customer, token, int32(u.cfg.JWT.Expiry)), nil
}
