package entity

import "time"

type CartRes struct {
	CartID     int32          `json:"cart_id"`
	CustomerID int32          `json:"customer_id"`
	Amount     int32          `json:"amount"`
	IsCheckout bool           `json:"is_checkout"`
	Items      []*CartItemRes `json:"items"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at,omitempty"`
}

type CartItemRes struct {
	CartItemID int32      `json:"cart_item_id"`
	CartID     int32      `json:"cart_id"`
	PartID     int32      `json:"part_id"`
	PartName   string     `json:"part_name"`
	SKU        string     `json:"sku"`
	Price      int32      `json:"price"`
	Quantity   int32      `json:"quantity"`
	Subtotal   int32      `json:"subtotal"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}
