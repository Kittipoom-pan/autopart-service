package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
	adminror "github.com/Kittipoom-pan/autopart-service/pkg/error"
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

func (r *adminRepository) GetAdminByID(ctx context.Context, id int) (*entitie.AdminRes, error) {
	admin, err := r.queries.GetAdminUser(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("admin_id", id).Msg("admin not found in database")
			return nil, adminror.NewNotFoundError("Admin")
		}
		r.logger.Error().Err(err).Int("admin_id", id).Msg("failed to get admin from database")
		return nil, adminror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	return entitie.MapDbAdminToAdminRes(admin), nil
}

func (r *adminRepository) CreateAdmin(ctx context.Context, admin *entitie.AdminReq) (int64, error) {
	params := entitie.MapAdminToAdminParam(admin, "SYSTEM")

	result, err := r.queries.CreateAdmin(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			r.logger.Warn().Err(mysqlErr).Str("admin", admin.Username).Msg("duplicate key error")
			return 0, adminror.NewAPIError(common.StatusConflict, "Username or email already exists")
		}
		r.logger.Error().Err(err).Msg("failed to create admin in database")
		return 0, adminror.NewAPIError(common.StatusError, "repository: failed to create admin")
	}

	adminID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve admin ID from database")
		return 0, adminror.NewAPIError(common.StatusError, "repository: failed to retrieve admin ID")
	}

	return adminID, nil
}

func (r *adminRepository) GetAllAdmins(ctx context.Context) ([]*entitie.AdminRes, error) {
	admins, err := r.queries.ListAdminUsers(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list admins from database")
		return nil, adminror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	var adminEntities []*entitie.AdminRes
	for _, admin := range admins {
		adminEntities = append(adminEntities, entitie.MapDbAdminsToAdminEntity(admin))
	}
	return adminEntities, nil
}

func (r *adminRepository) UpdateAdmin(ctx context.Context, id int, admin *entitie.AdminReq) error {
	params := entitie.MapUpdateAdminParams(id, admin, "SYSTEM")

	result, err := r.queries.UpdateAdmin(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			r.logger.Warn().Err(mysqlErr).Str("admin", admin.Username).Msg("duplicate key error")
			return adminror.NewAPIError(common.StatusConflict, "Username or email already exists")
		}
		r.logger.Error().Err(err).Int("admin_id", id).Msg("failed to update admin in database")
		return adminror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("admin_id", id).Msg("failed to get rows affected")
		return adminror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("admin_id", id).Msg("admin not found for update")
		return adminror.NewNotFoundError("Admin")
	}

	return nil
}

func (r *adminRepository) DeleteAdmin(ctx context.Context, id int) error {
	params := entitie.MapUpdateAdminIsActiveParams(id, false, "SYSTEM")

	result, err := r.queries.UpdateAdminIsActive(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int("admin_id", id).Msg("failed to delete admin in database")
		return adminror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("admin_id", id).Msg("failed to get rows affected")
		return adminror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("admin_id", id).Msg("admin not found for delete")
		return adminror.NewNotFoundError("Admin")
	}
	return nil
}

func (r *adminRepository) GetAdminByUsername(ctx context.Context, username string) (*entitie.Admin, error) {
	admin, err := r.queries.GetAdminByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Str("username", username).Msg("admin not found in database")
			return nil, adminror.NewNotFoundError("Admin")
		}
		r.logger.Error().Err(err).Str("username", username).Msg("failed to get admin from database")
		return nil, adminror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	return entitie.MapDbAdminToAdminEntity(admin), nil
}
