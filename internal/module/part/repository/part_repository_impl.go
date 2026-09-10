package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
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

func (r *partRepository) GetPartByID(ctx context.Context, id int) (*entity.PartRes, error) {
	part, err := r.queries.GetPartByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("part_id", id).Msg("part not found in database")
			return nil, parterror.NewNotFoundError("Part")
		}
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get part from database")
		return nil, parterror.NewAPIError(common.StatusError, "failed to get part")
	}

	return entity.MapDbPartToPartRes(part), nil
}

func (r *partRepository) CreatePart(ctx context.Context, part *entity.PartReq, createdBy *int) (int64, error) {
	params := entity.MapPartToPartParam(part, createdBy)
	result, err := r.queries.CreatePart(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				r.logger.Warn().Err(mysqlErr).Str("part", params.Name).Msg("duplicate key error")
				return 0, parterror.NewAPIError(common.StatusConflict, "Part already exists (duplicate key)")
			case 1452:
				r.logger.Warn().Err(mysqlErr).
					Int32("part_brand_id", params.PartBrandID).
					Int32("part_type_id", params.PartTypeID).
					Msg("foreign key constraint failed")
				return 0, parterror.NewAPIError(common.StatusBadRequest, "Invalid part_brand_id or part_type_id")
			}
		}
		r.logger.Error().Err(err).Str("sku", params.Sku).Msg("failed to create part in database")
		return 0, parterror.NewAPIError(common.StatusError, "failed to create part")
	}

	partID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve part ID from database")
		return 0, parterror.NewAPIError(common.StatusError, "failed to create part")
	}

	return partID, nil
}

func (r *partRepository) GetAllParts(ctx context.Context) ([]*entity.PartRes, error) {
	parts, err := r.queries.ListParts(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list parts from database")
		return nil, parterror.NewAPIError(common.StatusError, "failed to list parts")
	}

	partEntities := make([]*entity.PartRes, 0, len(parts))
	for _, part := range parts {
		partEntities = append(partEntities, entity.MapDbPartsToPartRes(part))
	}
	return partEntities, nil
}

func (r *partRepository) UpdatePart(ctx context.Context, id int, part *entity.PartReq, updatedBy int) error {
	params := entity.MapUpdatePartParams(id, part, updatedBy)
	result, err := r.queries.UpdatePartByID(ctx, params)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				r.logger.Warn().Err(mysqlErr).Str("part", params.Name).Msg("duplicate key error")
				return parterror.NewAPIError(common.StatusConflict, "Part already exists (duplicate key)")
			case 1452:
				r.logger.Warn().Err(mysqlErr).
					Int32("part_brand_id", params.PartBrandID).
					Int32("part_type_id", params.PartTypeID).
					Msg("foreign key constraint failed")
				return parterror.NewAPIError(common.StatusBadRequest, "Invalid part_brand_id or part_type_id")
			}
		}
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to update part in database")
		return parterror.NewAPIError(common.StatusError, "failed to update part")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get rows affected")
		return parterror.NewAPIError(common.StatusError, "failed to update part")
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("part_id", id).Msg("part not found for update")
		return parterror.NewNotFoundError("Part")
	}

	return nil
}

func (r *partRepository) DeletePart(ctx context.Context, id int, updatedBy int) error {
	params := entity.MapUpdatePartIsActiveParams(id, false, updatedBy)
	result, err := r.queries.DeletePartByID(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to delete part in database")
		return parterror.NewAPIError(common.StatusError, "failed to delete part")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get rows affected")
		return parterror.NewAPIError(common.StatusError, "failed to delete part")
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("part_id", id).Msg("part not found for delete")
		return parterror.NewNotFoundError("Part")
	}
	return nil
}
