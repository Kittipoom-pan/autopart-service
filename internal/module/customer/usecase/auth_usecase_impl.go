package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
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
		logger: log.With().Str("component", "auth_usecase").Logger(),
	}
}

func (u *authUsecase) Login(ctx context.Context, request *entitie.LoginRequest) (*entitie.LoginResponse, error) {
	customer, err := u.repo.GetCustomerByUsername(ctx, request.Username)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to get customer user")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(request.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			u.logger.Warn().Msg("Invalid password provided for user")
			return nil, customerror.NewAPIError(common.StatusUnauthorized, "Incorrect username or password")
		}

		u.logger.Error().Err(err).Msg("bcrypt comparison failed for unexpected reason")
		return nil, err
	}

	token, err := auth.GenerateToken(customer.ID, "", u.cfg)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to generate JWT token")
		return nil, err
	}

	return entitie.MapCustomerToLoginRes(customer, token, int32(u.cfg.JWT.Expiry)), nil
}
