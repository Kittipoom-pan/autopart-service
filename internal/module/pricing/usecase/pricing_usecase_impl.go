package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/pricing/entity"
	shippingEntity "github.com/Kittipoom-pan/autopart-service/internal/module/shipping/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/shipping/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
)

const DefaultUnknownUnitWeightG = 500

type pricingUsecase struct {
	shippingRepo repository.ShippingRateRepository
	logger       zerolog.Logger
}

func NewPricingUsecase(shippingRepo repository.ShippingRateRepository, logger zerolog.Logger) PricingUsecase {
	return &pricingUsecase{shippingRepo: shippingRepo, logger: logger}
}

func errAmountTooLarge() error {
	return customerror.InvalidRequestData(map[string]string{"items": "order amount too large"})
}

// CalculateCheckoutPricing calculates the subtotal, discount, shipping fee, and total amount.
// items ต้องประกอบจากข้อมูลฝั่ง server (DB) เท่านั้น ห้ามรับราคา/น้ำหนักจาก client
func (s *pricingUsecase) CalculateCheckoutPricing(
	ctx context.Context,
	items []entity.PricingItemInput,
	discountID *int32,
	dest *shippingEntity.ShippingDestination,
) (*entity.PricingResponse, error) {
	calculatedItems, subtotal, knownWeightG, unknownUnits, err := s.calculateItems(items)
	if err != nil {
		return nil, err
	}

	// discountAmount, err := s.calculateDiscount(ctx, subtotal, discountID)
	// if err != nil {
	// 	return nil, err
	// }

	//payableSubtotal := subtotal - discountAmount
	payableSubtotal := subtotal

	shippingFee, err := s.calculateShippingFee(ctx, dest, knownWeightG, unknownUnits, payableSubtotal)
	if err != nil {
		return nil, err
	}

	total := payableSubtotal + int64(shippingFee)
	if total > math.MaxInt32 {
		s.logger.Warn().Int64("total", total).Msg("checkout pricing total amount exceeds MaxInt32")
		return nil, errAmountTooLarge()
	}

	return &entity.PricingResponse{
		Items:    calculatedItems,
		Subtotal: int32(subtotal),
		//DiscountAmount: int32(discountAmount),
		ShippingFee: shippingFee,
		TotalAmount: int32(total),
	}, nil
}

// คำนวณสินค้าและยอดรวมเบื้องต้น, น้ำหนักรวมของชิ้นที่ทราบ (กรัม), จำนวนชิ้นที่ไม่ทราบน้ำหนัก
func (s *pricingUsecase) calculateItems(items []entity.PricingItemInput) ([]entity.CalculatedItem, int64, int64, int64, error) {
	if len(items) == 0 {
		return nil, 0, 0, 0, customerror.InvalidRequestData(map[string]string{"items": "cart is empty"})
	}

	var subtotal, knownWeightG, unknownUnits int64
	calculated := make([]entity.CalculatedItem, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, 0, 0, 0, customerror.InvalidRequestData(map[string]string{
				"items": fmt.Sprintf("part %d: quantity must be greater than 0", item.PartID),
			})
		}
		if item.Price < 0 {
			return nil, 0, 0, 0, customerror.InvalidRequestData(map[string]string{
				"items": fmt.Sprintf("part %d: price must not be negative", item.PartID),
			})
		}

		if item.WeightG == nil {
			unknownUnits += int64(item.Quantity) // นับเป็นจำนวนชิ้น ไม่ใช่จำนวนรายการ
		} else {
			if *item.WeightG < 0 {
				return nil, 0, 0, 0, customerror.InvalidRequestData(map[string]string{
					"items": fmt.Sprintf("part %d: weight must not be negative", item.PartID),
				})
			}
			knownWeightG += int64(*item.WeightG) * int64(item.Quantity)
		}

		itemSubtotal := int64(item.Price) * int64(item.Quantity)
		subtotal += itemSubtotal
		if subtotal > math.MaxInt32 || knownWeightG > math.MaxInt32 || unknownUnits > math.MaxInt32 {
			s.logger.Warn().
				Int64("subtotal", subtotal).
				Int64("known_weight_g", knownWeightG).
				Int64("unknown_units", unknownUnits).
				Msg("calculateItems values exceed MaxInt32")
			return nil, 0, 0, 0, errAmountTooLarge()
		}

		calculated[i] = entity.CalculatedItem{
			PartID:    item.PartID,
			Name:      item.Name,
			SKU:       item.SKU,
			Quantity:  item.Quantity,
			UnitPrice: item.Price,
			Subtotal:  int32(itemSubtotal),
			//DiscountAmount: 0,
			TotalPrice: int32(itemSubtotal),
		}
	}

	return calculated, subtotal, knownWeightG, unknownUnits, nil
}

// calculateDiscount ยังไม่รองรับ จึง reject ถ้ามี discountID ไม่เมินเฉย
// func (s *pricingUsecase) calculateDiscount(_ context.Context, subtotal int64, discountID *int32) (int64, error) {
// 	if discountID != nil {
// 		return 0, customerror.InvalidRequestData(map[string]string{"discount": "discount is not supported yet"})
// 	}

// 	var discount int64
// 	if discount > subtotal {
// 		discount = subtotal
// 	}
// 	return discount, nil
// }

func (s *pricingUsecase) calculateShippingFee(
	ctx context.Context,
	dest *shippingEntity.ShippingDestination,
	knownWeightG int64,
	unknownUnits int64,
	payableSubtotal int64,
) (int32, error) {
	var provinceID *int32
	if dest != nil {
		provinceID = dest.ProvinceID
	}

	// ชิ้นที่ไม่ทราบน้ำหนักนับเป็น unknownUnitWeightG ต่อชิ้น
	// (int64 × int64 ไม่ overflow เพราะแต่ละค่าไม่เกิน MaxInt32)
	//effectiveWeightG := knownWeightG + unknownUnits*int64(s.unknownUnitWeightG) // TODO: ปรับให้เป็น config ต่อไป
	effectiveWeightG := knownWeightG + unknownUnits*int64(DefaultUnknownUnitWeightG)
	if effectiveWeightG > math.MaxInt32 {
		s.logger.Warn().Int64("effective_weight_g", effectiveWeightG).Msg("effective weight exceeds MaxInt32")
		return 0, errAmountTooLarge()
	}

	rate, err := s.shippingRepo.FindShippingRate(ctx, provinceID, int32(effectiveWeightG))
	if err != nil {
		if errors.Is(err, repository.ErrShippingRateNotFound) {
			return 0, customerror.InvalidRequestData(map[string]string{"shipping": "shipping is not available for this destination"})
		}
		return 0, helper.MapDBErrorToAPIError(err, "Shipping rate")
	}

	// ถ้ายอดถึงเกณฑ์แล้วค่าส่งจะกลายเป็น 0 บาท
	if rate.FreeShippingMinSubtotal != nil && payableSubtotal >= int64(*rate.FreeShippingMinSubtotal) {
		s.logger.Info().Int64("payable_subtotal", payableSubtotal).Msg("Free shipping applied")
		return 0, nil
	}
	return rate.Fee, nil
}
