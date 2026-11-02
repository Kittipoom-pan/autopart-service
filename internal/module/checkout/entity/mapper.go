package entity

import (
	"fmt"

	addressEntity "github.com/Kittipoom-pan/autopart-service/internal/module/address/entity"
	pricingEntity "github.com/Kittipoom-pan/autopart-service/internal/module/pricing/entity"
)

func FormatMoney(amount int32) string {
	return fmt.Sprintf("%.2f", float64(amount))
}

func MapPricingResultToCheckoutPreviewRes(pricing *pricingEntity.PricingResponse, address *addressEntity.AddressRes) *CheckoutPreviewRes {
	if pricing == nil {
		return nil
	}

	items := make([]CheckoutPreviewItem, len(pricing.Items))
	for i, item := range pricing.Items {
		items[i] = CheckoutPreviewItem{
			PartID:    item.PartID,
			Name:      item.Name,
			SKU:       item.SKU,
			Quantity:  item.Quantity,
			UnitPrice: FormatMoney(item.UnitPrice),
			Subtotal:  FormatMoney(item.Subtotal),
			//DiscountAmount: FormatMoney(item.DiscountAmount),
			TotalPrice: FormatMoney(item.TotalPrice),
		}
	}

	return &CheckoutPreviewRes{
		Address:  address,
		Items:    items,
		Subtotal: FormatMoney(pricing.Subtotal),
		//DiscountAmount: FormatMoney(pricing.DiscountAmount),
		ShippingFee: FormatMoney(pricing.ShippingFee),
		TotalAmount: FormatMoney(pricing.TotalAmount),
	}
}
