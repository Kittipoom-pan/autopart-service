package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/image/entity"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type imageRepository struct {
	queries *db.Queries
	logger  zerolog.Logger
}

func NewImageRepository(conn *sql.DB, logger zerolog.Logger) ImageRepository {
	return &imageRepository{
		queries: db.New(conn),
		logger:  logger,
	}
}

func (r *imageRepository) getQueries(ctx context.Context) *db.Queries {
	if tx, ok := database.GetTx(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *imageRepository) ListImagesByReference(ctx context.Context, referenceType string, referenceID int32) ([]entity.Image, error) {
	rows, err := r.getQueries(ctx).ListImagesByReference(ctx, db.ListImagesByReferenceParams{
		ReferenceType: referenceType,
		ReferenceID:   referenceID,
	})
	if err != nil {
		r.logger.Error().Err(err).
			Str("reference_type", referenceType).
			Int32("reference_id", referenceID).
			Msg("failed to list images by reference")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	images := make([]entity.Image, 0, len(rows))
	for _, row := range rows {
		images = append(images, entity.MapDbImageToImage(row))
	}

	return images, nil
}

func (r *imageRepository) CreateImage(ctx context.Context, input CreateImageInput) (int64, error) {
	result, err := r.getQueries(ctx).AddImage(ctx, db.AddImageParams{
		ReferenceID:      input.ReferenceID,
		ReferenceType:    input.ReferenceType,
		FileName:         input.FileName,
		ImageUrl:         input.ImageURL,
		AltText:          input.AltText,
		IsPrimary:        input.IsPrimary,
		SortImage:        input.SortImage,
		CreatedAt:        utils.NowUTC(),
		CreatedByAdminID: utils.IntNullToNullInt32(input.CreatedByAdminID),
	})
	if err != nil {
		dbErr := dberr.HandleMySQLError(err)
		r.logger.Error().Err(err).
			Str("reference_type", input.ReferenceType).
			Int32("reference_id", input.ReferenceID).
			Msg("failed to create image in database")
		return 0, dbErr
	}

	imageID, err := result.LastInsertId()
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to retrieve image id from database")
		return 0, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return imageID, nil
}

func (r *imageRepository) DeleteImage(ctx context.Context, imageID int64) error {
	if _, err := r.getQueries(ctx).DeleteImage(ctx, imageID); err != nil {
		r.logger.Error().Err(err).Int64("image_id", imageID).Msg("failed to delete image from database")
		return dberr.HandleMySQLError(err)
	}

	return nil
}
