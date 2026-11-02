package entity

import (
	addressEntity "github.com/Kittipoom-pan/autopart-service/internal/module/address/entity"
)

type CheckoutPreviewRes struct {
	Address        *addressEntity.AddressRes `json:"address,omitempty"`
	Items          []CheckoutPreviewItem     `json:"items"`
	Subtotal       string                    `json:"subtotal"`
	DiscountAmount string                    `json:"discount_amount"`
	ShippingFee    string                    `json:"shipping_fee"`
	TotalAmount    string                    `json:"total_amount"`
}

type CheckoutPreviewItem struct {
	PartID         int32  `json:"part_id"`
	Name           string `json:"name"`
	SKU            string `json:"sku"`
	Quantity       int32  `json:"quantity"`
	UnitPrice      string `json:"unit_price"`
	Subtotal       string `json:"subtotal"`
	DiscountAmount string `json:"discount_amount"`
	TotalPrice     string `json:"total_price"`
}
