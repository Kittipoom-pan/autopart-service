package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/auth"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/repository"
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

func (u *adminUsecase) GetAdminByID(ctx context.Context, id int) (*entitie.Admin, error) {
	u.logger.Info().Int("admin_id", id).Msg("GetAdminByID started")

	admin, err := u.repo.GetAdminByID(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Int("admin_id", id).Msg("Failed to get admin from repository")
		return nil, err
	}

	u.logger.Info().Int("admin_id", id).Msg("GetAdminByID completed successfully")
	return admin, nil
}

func (u *adminUsecase) GetAllAdmins(ctx context.Context) ([]*entitie.Admin, error) {
	admins, err := u.repo.GetAllAdmins(ctx)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to get all admins from repository")
		return nil, err
	}
	return admins, nil
}

func (u *adminUsecase) CreateAdmin(ctx context.Context, admin *entitie.AdminReq) (int64, error) {
	u.logger.Info().Msg("CreateAdmin started")

	hashedPassword, err := auth.HashPassword(admin.Password)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to hash password")
		return 0, err
	}

	admin.Password = hashedPassword
	adminID, err := u.repo.CreateAdmin(ctx, admin)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to create admin")
		return 0, err
	}
	u.logger.Info().Int64("admin_id", adminID).Msg("Admin created successfully")
	return adminID, nil
}

func (u *adminUsecase) UpdateAdmin(ctx context.Context, id int, user *entitie.AdminReq) error {
	u.logger.Info().Int("admin_id", id).Msg("UpdateAdmin started")

	hashedPassword, err := auth.HashPassword(user.Password)

	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to hash password")
		return err
	}

	user.Password = hashedPassword

	err = u.repo.UpdateAdmin(ctx, id, user)
	if err != nil {
		u.logger.Error().Err(err).Int("admin_id", id).Msg("Failed to update admin in repository")
		return err
	}

	u.logger.Info().Int("admin_id", id).Msg("UpdateAdmin completed successfully")
	return nil
}

func (u *adminUsecase) DeleteAdmin(ctx context.Context, id int) error {
	u.logger.Info().Int("admin_id", id).Msg("DeleteAdmin started")

	err := u.repo.DeleteAdmin(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Int("admin_id", id).Msg("Failed to delete admin in repository")
		return err
	}

	u.logger.Info().Int("admin_id", id).Msg("DeleteAdmin completed successfully")
	return nil
}
