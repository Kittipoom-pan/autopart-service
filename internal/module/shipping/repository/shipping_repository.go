package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/shipping/entity"
)

type ShippingRateRepository interface {
	FindShippingRate(ctx context.Context, provinceID *int32, weightG int32) (*entity.ShippingRate, error)
}
