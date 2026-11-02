package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/pricing/entity"
	shippingEntity "github.com/Kittipoom-pan/autopart-service/internal/module/shipping/entity"
)

type PricingUsecase interface {
	CalculateCheckoutPricing(
		ctx context.Context,
		items []entity.PricingItemInput,
		discountID *int32,
		dest *shippingEntity.ShippingDestination,
	) (*entity.PricingResponse, error)
}
