package usecase

import (
	"context"
	"fmt"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	addressRepository "github.com/Kittipoom-pan/autopart-service/internal/module/address/repository"
	cartEntity "github.com/Kittipoom-pan/autopart-service/internal/module/cart/entity"
	cartRepository "github.com/Kittipoom-pan/autopart-service/internal/module/cart/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout/usecase/validation"
	partRepository "github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	pricingEnitity "github.com/Kittipoom-pan/autopart-service/internal/module/pricing/entity"
	pricingUsecase "github.com/Kittipoom-pan/autopart-service/internal/module/pricing/usecase"
	shippingEntity "github.com/Kittipoom-pan/autopart-service/internal/module/shipping/entity"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
)

type checkoutUsecase struct {
	cartRepo       cartRepository.CartRepository
	partRepo       partRepository.PartRepository
	addressRepo    addressRepository.AddressRepository
	pricingUsecase pricingUsecase.PricingUsecase
	logger         zerolog.Logger
}

func NewCheckoutUsecase(
	cartRepo cartRepository.CartRepository,
	partRepo partRepository.PartRepository,
	addressRepo addressRepository.AddressRepository,
	pricingUsecase pricingUsecase.PricingUsecase,
	logger zerolog.Logger,
) CheckoutUsecase {
	return &checkoutUsecase{
		cartRepo:       cartRepo,
		partRepo:       partRepo,
		addressRepo:    addressRepo,
		pricingUsecase: pricingUsecase,
		logger:         logger,
	}
}

func (u *checkoutUsecase) PreviewCheckout(ctx context.Context, customerID int32, req *entity.CheckoutPreviewReq) (*entity.CheckoutPreviewRes, error) {
	u.logger.Debug().
		Int32("customer_id", customerID).
		Int32("cart_id", req.CartID).
		Int32("address_id", req.AddressID).
		Msg("PreviewCheckout")

	if err := validation.ValidateCheckoutPreviewRequest(req); err != nil {
		u.logger.Warn().Err(err).Msg("checkout preview validation failed")
		return nil, err
	}

	// check if the address belongs to the customer
	address, err := u.addressRepo.GetAddressByIDAndCustomer(ctx, req.AddressID, customerID)
	if err != nil {
		u.logger.Warn().Err(err).Int32("address_id", req.AddressID).Int32("customer_id", customerID).Msg("address not found or not owned by customer")
		return nil, helper.MapDBErrorToAPIError(err, "Address")
	}

	cart, err := u.cartRepo.GetActiveCartByCustomerID(ctx, customerID)
	if err != nil {
		u.logger.Warn().Err(err).Int32("customer_id", customerID).Msg("active cart not found")
		return nil, helper.MapDBErrorToAPIError(err, "Cart")
	}

	if cart.CartID != req.CartID {
		u.logger.Warn().
			Int32("req_cart_id", req.CartID).
			Int32("active_cart_id", cart.CartID).
			Msg("cart id mismatch")
		return nil, customerror.NewAPIError(common.StatusNotFound, "cart not found")
	}

	if len(cart.Items) == 0 {
		u.logger.Warn().Int32("cart_id", cart.CartID).Msg("cart is empty")
		return nil, customerror.NewAPIError(common.StatusBadRequest, "cart is empty")
	}

	pricingItems, err := u.buildPricingItems(ctx, cart.Items)
	if err != nil {
		return nil, err
	}

	dest := &shippingEntity.ShippingDestination{ProvinceID: address.ProvinceID}

	pricingResult, err := u.pricingUsecase.CalculateCheckoutPricing(ctx, pricingItems, req.DiscountID, dest)
	if err != nil {
		u.logger.Warn().Err(err).Msg("failed to calculate pricing")
		return nil, err
	}

	return entity.MapPricingResultToCheckoutPreviewRes(pricingResult, address), nil
}

func (u *checkoutUsecase) buildPricingItems(ctx context.Context, cartItems []*cartEntity.CartItemRes) ([]pricingEnitity.PricingItemInput, error) {
	ids := make([]int32, 0, len(cartItems))
	for _, item := range cartItems {
		ids = append(ids, item.PartID)
	}

	parts, err := u.partRepo.GetPartsByIDs(ctx, ids)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Part")
	}

	items := make([]pricingEnitity.PricingItemInput, 0, len(cartItems))
	for _, item := range cartItems {
		part, ok := parts[item.PartID]
		if !ok {
			u.logger.Warn().Int32("part_id", item.PartID).Msg("part not found or inactive")
			return nil, customerror.NewAPIError(common.StatusBadRequest,
				fmt.Sprintf("part %d is no longer available", item.PartID))
		}

		if item.Quantity > part.Quantity {
			u.logger.Warn().
				Int32("part_id", item.PartID).
				Str("part_name", part.Name).
				Int32("available_stock", part.Quantity).
				Int32("requested_quantity", item.Quantity).
				Msg("insufficient stock")
			return nil, customerror.NewAPIError(
				common.StatusBadRequest,
				fmt.Sprintf("insufficient stock for part '%s' (available: %d, requested: %d)", part.Name, part.Quantity, item.Quantity),
			)
		}

		if part.Price < 0 || (part.WeightG != nil && *part.WeightG < 0) {
			u.logger.Error().Int32("part_id", part.PartID).Msg("part has invalid price or weight")
			return nil, customerror.NewAPIError(common.StatusInternalServer, "invalid part data")
		}

		if part.WeightG == nil {
			u.logger.Warn().Int32("part_id", part.PartID).Msg("part has no weight, fallback weight will be used for shipping")
		}

		items = append(items, pricingEnitity.PricingItemInput{
			PartID:   part.PartID,
			Name:     part.Name,
			SKU:      part.SKU,
			Quantity: item.Quantity,
			Price:    part.Price,
			WeightG:  part.WeightG,
		})
	}

	return items, nil
}
