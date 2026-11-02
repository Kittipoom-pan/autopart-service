package repository

import (
	"context"
	"database/sql"
	"errors"

	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/shipping/entity"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

var ErrShippingRateNotFound = errors.New("shipping rate not found")


type shippingRateRepository struct {
	queries *db.Queries
	logger  zerolog.Logger
}

func NewShippingRateRepository(conn *sql.DB, logger zerolog.Logger) ShippingRateRepository {
	return &shippingRateRepository{queries: db.New(conn), logger: logger}
}

func (r *shippingRateRepository) FindShippingRate(ctx context.Context, provinceID *int32, weightG int32) (*entity.ShippingRate, error) {
	w := int64(weightG)

	row, err := r.queries.FindShippingRate(ctx, db.FindShippingRateParams{
		ProvinceID: utils.Int32PtrToNullInt32(provinceID),
		WeightG:    w,
		WeightG_2:  w,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrShippingRateNotFound
		}
		r.logger.Error().Err(err).Int32("weight_g", weightG).Msg("failed to find shipping rate")
		return nil, dberr.NewDBError(dberr.ErrDatabaseInternal, err.Error())
	}

	return &entity.ShippingRate{
		Fee:                     row.Fee,
		FreeShippingMinSubtotal: utils.NullInt32ToInt32Ptr(row.FreeShippingMinSubtotal),
	}, nil
}
