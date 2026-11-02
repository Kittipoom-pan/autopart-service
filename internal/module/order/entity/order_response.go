package entity

import (
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
)

type OrderRes struct {
	OrderID             int64                    `json:"order_id"`
	OrderNumber         string                   `json:"order_number"`
	CustomerID          *int32                   `json:"customer_id,omitempty"`
	Channel             db.OrderChannel          `json:"channel"`
	PlacedByAdminUserID *int32                   `json:"placed_by_admin_user_id,omitempty"`
	CartID              *int32                   `json:"cart_id,omitempty"`
	Subtotal            int32                    `json:"subtotal"`
	DiscountAmount      int32                    `json:"discount_amount"`
	ShippingFee         int32                    `json:"shipping_fee"`
	TotalAmount         int32                    `json:"total_amount"`
	DiscountID          *int32                   `json:"discount_id,omitempty"`
	Status              db.OrderStatus           `json:"status"`
	Remark              *string                  `json:"remark,omitempty"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           *time.Time               `json:"updated_at,omitempty"`
	Items               []OrderItemRes           `json:"items,omitempty"`
	ShippingAddress     *OrderShippingAddressRes `json:"shipping_address,omitempty"`
}

type OrderItemRes struct {
	OrderItemID int64 `json:"order_item_id"`
	OrderID     int64 `json:"order_id"`
	PartID      int32 `json:"part_id"`
	Quantity    int32 `json:"quantity"`
	Price       int32 `json:"price"`
}

type OrderShippingAddressRes struct {
	OrderID       int64  `json:"order_id"`
	RecipientName string `json:"recipient_name"`
	PhoneNumber   string `json:"phone_number"`
	AddressLine   string `json:"address_line"`
	Subdistrict   string `json:"subdistrict"`
	District      string `json:"district"`
	Province      string `json:"province"`
	ZipCode       string `json:"zip_code"`
}
