package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
	adminerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type adminRepository struct {
	queries *db.Queries
	logger  zerolog.Logger
}

func NewAdminRepository(queries *db.Queries) AdminRepository {
	return &adminRepository{
		queries: queries,
		logger:  log.With().Str("component", "admin_repository").Logger(),
	}
}

func (r *adminRepository) GetAdminByID(ctx context.Context, id int) (*entity.AdminRes, error) {
	admin, err := r.queries.GetAdminUser(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("admin_id", id).Msg("admin not found in database")
			return nil, adminerror.NewNotFoundError("Admin")
		}
		r.logger.Error().Err(err).Int("admin_id", id).Msg("failed to get admin from database")
		return nil, adminerror.NewAPIError(common.StatusError, "failed to get admin")
	}

	return entity.MapDbAdminToAdminRes(admin), nil
}

func (r *adminRepository) CreateAdmin(ctx context.Context, admin *entity.AdminReq, createdBy *int) (int64, error) {
	params := entity.MapAdminToAdminParam(admin, createdBy)
	result, err := r.queries.CreateAdmin(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			r.logger.Warn().Err(mysqlErr).Str("admin", params.Username).Msg("duplicate key error")
			return 0, adminerror.NewAPIError(common.StatusConflict, "Username or email already exists")
		}
		r.logger.Error().Err(err).Msg("failed to create admin in database")
		return 0, adminerror.NewAPIError(common.StatusError, "failed to create admin")
	}

	adminID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve admin ID from database")
		return 0, adminerror.NewAPIError(common.StatusError, "failed to create admin")
	}

	return adminID, nil
}

func (r *adminRepository) GetAllAdmins(ctx context.Context) ([]*entity.AdminRes, error) {
	admins, err := r.queries.ListAdminUsers(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list admins from database")
		return nil, adminerror.NewAPIError(common.StatusError, "failed to list admins")
	}

	adminEntities := make([]*entity.AdminRes, 0, len(admins))
	for _, admin := range admins {
		adminEntities = append(adminEntities, entity.MapDbAdminsToAdminEntity(admin))
	}
	return adminEntities, nil
}

func (r *adminRepository) UpdateAdmin(ctx context.Context, adminID int, admin *entity.AdminReq, updatedBy int) error {
	params := entity.MapUpdateAdminParams(adminID, admin, &updatedBy)
	result, err := r.queries.UpdateAdmin(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			r.logger.Warn().Err(mysqlErr).Str("admin", params.Username).Msg("duplicate key error")
			return adminerror.NewAPIError(common.StatusConflict, "Username or email already exists")
		}
		r.logger.Error().Err(err).Int("admin_id", adminID).Msg("failed to update admin in database")
		return adminerror.NewAPIError(common.StatusError, "failed to update admin")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("admin_id", adminID).Msg("failed to get rows affected")
		return adminerror.NewAPIError(common.StatusError, "failed to update admin")
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("admin_id", adminID).Msg("admin not found for update")
		return adminerror.NewNotFoundError("Admin")
	}

	return nil
}

func (r *adminRepository) DeleteAdmin(ctx context.Context, adminID int, updatedBy int) error {
	params := entity.MapUpdateAdminIsActiveParams(adminID, false, &updatedBy)
	result, err := r.queries.UpdateAdminIsActive(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int("admin_id", adminID).Msg("failed to delete admin in database")
		return adminerror.NewAPIError(common.StatusError, "failed to delete admin")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("admin_id", adminID).Msg("failed to get rows affected")
		return adminerror.NewAPIError(common.StatusError, "failed to delete admin")
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("admin_id", adminID).Msg("admin not found for delete")
		return adminerror.NewNotFoundError("Admin")
	}
	return nil
}

func (r *adminRepository) GetAdminByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	admin, err := r.queries.GetAdminByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Str("username", username).Msg("admin not found in database")
			return nil, adminerror.NewNotFoundError("Admin")
		}
		r.logger.Error().Err(err).Str("username", username).Msg("failed to get admin from database")
		return nil, adminerror.NewAPIError(common.StatusError, "failed to get admin")
	}

	return entity.MapDbAdminToAdminEntity(admin), nil
}
