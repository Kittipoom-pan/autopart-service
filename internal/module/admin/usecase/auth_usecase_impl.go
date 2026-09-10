package usecase

import (
	"context"
	"errors"

	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		logger: log.With().Str("component", "admin_auth_usecase").Logger(),
	}
}

func (u *authUsecase) Login(ctx context.Context, request *entity.LoginRequest) (*entity.LoginResponse, error) {
	invalidCreds := customerror.NewAPIError(common.StatusUnauthorized, "Incorrect username or password")

	admin, err := u.repo.GetAdminByUsername(ctx, request.Username)
	if err != nil {
		var notFound *customerror.NotFoundError
		if errors.As(err, &notFound) {
			u.logger.Warn().Str("username", request.Username).Msg("login failed: admin not found")
			return nil, invalidCreds
		}
		u.logger.Error().Err(err).Msg("failed to get admin for login")
		return nil, err
	}

	if err := auth.VerifyPassword(admin.Password, request.Password); err != nil {
		u.logger.Warn().Str("username", request.Username).Msg("login failed: invalid password")
		return nil, invalidCreds
	}

	token, err := auth.GenerateToken(admin.ID, admin.Role, u.cfg)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to generate JWT token")
		return nil, customerror.NewAPIError(common.StatusError, "failed to generate token")
	}

	return entity.MapAdminToLoginRes(admin, token, int32(u.cfg.JWT.Expiry)), nil
}
