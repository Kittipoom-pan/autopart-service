package repository

import (
	"context"
	"database/sql"

	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
	"github.com/rs/zerolog"
)

type partRepository struct {
	queries *db.Queries
	logger  zerolog.Logger
}

func NewPartRepository(queries *db.Queries, logger zerolog.Logger) PartRepository {
	return &partRepository{
		queries: queries,
		logger:  logger,
	}
}

func (r *partRepository) GetPartByID(ctx context.Context, id int) (*entity.PartRes, error) {
	part, err := r.queries.GetPartByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Int("part_id", id).Msg("part not found in database")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "part not found")
		}
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get part from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbPartToPartRes(part), nil
}

func (r *partRepository) GetPartsByIDs(ctx context.Context, ids []int32) (map[int32]*entity.PartRes, error) {
	result := make(map[int32]*entity.PartRes, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	// ตัด id ซ้ำ (cart ไม่ควรมี part ซ้ำ แต่กันไว้)
	seen := make(map[int32]struct{}, len(ids))
	unique := make([]int32, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			unique = append(unique, id)
		}
	}

	rows, err := r.queries.GetPartsByIDs(ctx, unique)
	if err != nil {
		r.logger.Error().Err(err).Int("count", len(unique)).Msg("failed to get parts by ids")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	for _, row := range rows {
		result[row.PartID] = entity.MapDbPartToPartRes(db.GetPartByIDRow(row))
	}
	return result, nil
}

func (r *partRepository) GetPartBySKU(ctx context.Context, sku string) (*entity.PartRes, error) {
	part, err := r.queries.GetPartBySKU(ctx, sku)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn().Str("sku", sku).Msg("part not found in database")
			return nil, dberr.NewDBError(dberr.ErrRecordNotFound, "part not found")
		}
		r.logger.Error().Err(err).Str("sku", sku).Msg("failed to get part from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return entity.MapDbPartBySKUToPartRes(part), nil
}

func (r *partRepository) CreatePart(ctx context.Context, part *entity.PartReq, createdBy *int) (int64, error) {
	params := entity.MapPartToPartParam(part, createdBy)
	result, err := r.queries.CreatePart(ctx, params)
	if err != nil {
		dbErr := dberr.HandleMySQLError(err)
		if dberr.IsDuplicateKey(dbErr) {
			r.logger.Warn().Err(dbErr).Str("part", params.Name).Msg("duplicate key error")
			return 0, dbErr
		}
		if dberr.IsForeignKeyError(dbErr) {
			r.logger.Warn().Err(dbErr).
				Int32("part_brand_id", params.PartBrandID).
				Int32("part_type_id", params.PartTypeID).
				Msg("foreign key constraint failed")
			return 0, dbErr
		}
		r.logger.Error().Err(err).Str("sku", params.Sku).Msg("original database error")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	partID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve part ID from database")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return partID, nil
}

func (r *partRepository) GetAllParts(ctx context.Context) ([]*entity.PartRes, error) {
	parts, err := r.queries.ListParts(ctx)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list parts from database")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
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
		dbErr := dberr.HandleMySQLError(err)
		if dberr.IsDuplicateKey(dbErr) {
			r.logger.Warn().Err(dbErr).Str("part", params.Name).Msg("duplicate key error")
			return dbErr
		}
		if dberr.IsForeignKeyError(dbErr) {
			r.logger.Warn().Err(dbErr).
				Int32("part_brand_id", params.PartBrandID).
				Int32("part_type_id", params.PartTypeID).
				Msg("foreign key constraint failed")
			return dbErr
		}
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to update part in database")
		return dbErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get rows affected")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("part_id", id).Msg("part not found for update")
		return dberr.NewDBError(dberr.ErrRecordNotFound, "part not found")
	}

	return nil
}

func (r *partRepository) DeletePart(ctx context.Context, id int, updatedBy int) error {
	params := entity.MapUpdatePartIsActiveParams(id, false, updatedBy)
	result, err := r.queries.DeletePartByID(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to delete part in database")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error().Err(err).Int("part_id", id).Msg("failed to get rows affected")
		return dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	if rowsAffected == 0 {
		r.logger.Warn().Int("part_id", id).Msg("part not found for delete")
		return dberr.NewDBError(dberr.ErrRecordNotFound, "part not found")
	}
	return nil
}
