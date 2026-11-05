package entity

type CartRes struct {
	CartID   int32          `json:"cart_id"`
	Subtotal int32          `json:"subtotal"`
	Items    []*CartItemRes `json:"items"`
}

type CartItemRes struct {
	CartItemID int32  `json:"cart_item_id"`
	PartID     int32  `json:"part_id"`
	PartName   string `json:"part_name"`
	SKU        string `json:"sku"`
	Price      int32  `json:"price"`
	Quantity   int32  `json:"quantity"`
	Subtotal   int32  `json:"subtotal"` // Price * Quantity
}
