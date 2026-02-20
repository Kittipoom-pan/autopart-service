package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entitie"
	parterror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type partRepository struct {
	queries *db.Queries
	logger  zerolog.Logger
}

func NewPartRepository(queries *db.Queries) PartRepository {
	return &partRepository{
		queries: queries,
		logger:  log.With().Str("component", "part_repository").Logger(),
	}
}

func (r *partRepository) GetPartByID(ctx context.Context, id int) (*entitie.PartRes, error) {
	part, err := r.queries.GetPartByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("part_id", id).Msg("part not found in database")
			return nil, parterror.NewNotFoundError("Part")
		}
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get part from database")
		return nil, parterror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	return entitie.MapDbPartToPartRes(part), nil
}

func (r *partRepository) CreatePart(ctx context.Context, part *entitie.PartReq) (int64, error) {
	params := entitie.MapPartToPartParam(part, "SYSTEM")

	result, err := r.queries.CreatePart(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			r.logger.Warn().Err(mysqlErr).Str("part", part.Name).Msg("duplicate key error")
			return 0, parterror.NewAPIError(common.StatusConflict, "Part already exists (duplicate key)")
		}
		r.logger.Error().Err(err).Msg("failed to create part in database")
		return 0, parterror.NewAPIError(common.StatusError, "repository: failed to create part")
	}

	partID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve part ID from database")
		return 0, parterror.NewAPIError(common.StatusError, "repository: failed to retrieve part ID")
	}

	return partID, nil
}

func (r *partRepository) GetAllParts(ctx context.Context) ([]*entitie.PartRes, error) {
	parts, err := r.queries.ListParts(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list parts from database")
		return nil, parterror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	var partEntities []*entitie.PartRes
	for _, part := range parts {
		partEntities = append(partEntities, entitie.MapDbPartsToPartRes(part))
	}
	return partEntities, nil
}

func (r *partRepository) UpdatePart(ctx context.Context, id int, part *entitie.PartReq) error {
	params := entitie.MapUpdatePartParams(id, part, "SYSTEM")

	result, err := r.queries.UpdatePartByID(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			r.logger.Warn().Err(mysqlErr).Str("part", part.Name).Msg("duplicate key error")
			return parterror.NewAPIError(common.StatusConflict, "Part already exists (duplicate key)")
		}
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to update part in database")
		return parterror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get rows affected")
		return parterror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("part_id", id).Msg("part not found for update")
		return parterror.NewNotFoundError("Part")
	}

	return nil
}

func (r *partRepository) DeletePart(ctx context.Context, id int) error {
	params := entitie.MapUpdatePartIsActiveParams(id, false, "SYSTEM")

	result, err := r.queries.DeletePartByID(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to delete part in database")
		return parterror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get rows affected")
		return parterror.NewAPIError(common.StatusError, "Database error: "+err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("part_id", id).Msg("part not found for delete")
		return parterror.NewNotFoundError("Part")
	}
	return nil
}
