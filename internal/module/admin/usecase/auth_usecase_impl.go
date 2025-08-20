package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo   repository.AdminRepository
	logger zerolog.Logger
	cfg    *config.Config
}

func NewAuthUsecase(repo repository.AdminRepository, cfg *config.Config) AuthUsecase {
	return &authUsecase{
		repo:   repo,
		cfg:    cfg,
		logger: log.With().Str("component", "auth_usecase").Logger(),
	}
}

func (u *authUsecase) Login(ctx context.Context, request *entitie.LoginRequest) (*entitie.LoginResponse, error) {
	admin, err := u.repo.GetAdminByUsername(ctx, request.Username)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to get admin user")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			u.logger.Warn().Msg("Invalid password provided for user")
			return nil, customerror.NewAPIError(common.StatusUnauthorized, "Incorrect username or password")
		}

		u.logger.Error().Err(err).Msg("bcrypt comparison failed for unexpected reason")
		return nil, err
	}

	token, err := auth.GenerateToken(admin.ID, admin.Role, u.cfg)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to generate JWT token")
		return nil, err
	}

	return entitie.MapAdminToLoginRes(admin, token, int32(u.cfg.JWT.Expiry)), nil
}
