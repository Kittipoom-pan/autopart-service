package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/usecase/validation"
	adminerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type adminUsecase struct {
	repo   repository.AdminRepository
	logger zerolog.Logger
}

func NewAdminUsecase(repo repository.AdminRepository) AdminUsecase {
	return &adminUsecase{
		repo:   repo,
		logger: log.With().Str("component", "admin_usecase").Logger(),
	}
}

func (u *adminUsecase) GetAdminByID(ctx context.Context, id int) (*entity.AdminRes, error) {
	u.logger.Debug().Int("admin_id", id).Msg("GetAdminByID")
	return u.repo.GetAdminByID(ctx, id)
}

func (u *adminUsecase) GetAllAdmins(ctx context.Context) ([]*entity.AdminRes, error) {
	u.logger.Debug().Msg("GetAllAdmins")
	return u.repo.GetAllAdmins(ctx)
}

func (u *adminUsecase) CreateAdmin(ctx context.Context, admin *entity.AdminReq, userID *int) (int64, error) {
	if err := validation.ValidateAdminRequest(admin, false); err != nil {
		u.logger.Warn().Err(err).Msg("admin create validation failed")
		return 0, err
	}

	hashedPassword, err := auth.HashPassword(admin.Password)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to hash password")
		return 0, adminerror.NewAPIError(common.StatusError, "failed to process password")
	}

	req := *admin
	req.Password = hashedPassword

	adminID, err := u.repo.CreateAdmin(ctx, &req, userID)
	if err != nil {
		return 0, err
	}

	u.logger.Info().Int64("admin_id", adminID).Msg("admin created successfully")
	return adminID, nil
}

func (u *adminUsecase) UpdateAdmin(ctx context.Context, adminID int, userID int, admin *entity.AdminReq) error {
	if err := validation.ValidateAdminRequest(admin, true); err != nil {
		u.logger.Warn().Err(err).Msg("admin update validation failed")
		return err
	}

	hashedPassword, err := auth.HashPassword(admin.Password)
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to hash password")
		return adminerror.NewAPIError(common.StatusError, "failed to process password")
	}

	req := *admin
	req.Password = hashedPassword

	if err := u.repo.UpdateAdmin(ctx, adminID, &req, userID); err != nil {
		return err
	}

	u.logger.Info().Int("admin_id", adminID).Msg("admin updated successfully")
	return nil
}

func (u *adminUsecase) DeleteAdmin(ctx context.Context, adminID int, userID int) error {
	if err := u.repo.DeleteAdmin(ctx, adminID, userID); err != nil {
		return err
	}

	u.logger.Info().Int("admin_id", adminID).Msg("admin deleted successfully")
	return nil
}
