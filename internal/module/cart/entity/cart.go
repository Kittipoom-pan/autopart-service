package entity

import (
	"time"
)

type Cart struct {
	CartID     int32
	CustomerID int32
	Amount     int32
	IsCheckout bool
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

type CartItem struct {
	CartItemID int32
	CartID     int32
	PartID     int32
	Quantity   int32
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
