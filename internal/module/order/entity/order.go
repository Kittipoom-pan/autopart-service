package entity

import (
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
)

type Order struct {
	OrderID             int64
	OrderNumber         string
	CustomerID          *int32
	Channel             db.OrderChannel
	PlacedByAdminUserID *int32
	CartID              *int32
	Subtotal            int32
	DiscountAmount      int32
	ShippingFee         int32
	TotalAmount         int32
	DiscountID          *int32
	Status              db.OrderStatus
	Remark              *string
	CreatedAt           time.Time
	UpdatedAt           *time.Time
}

type OrderItem struct {
	OrderItemID int64
	OrderID     int64
	PartID      int32
	Quantity    int32
	Price       int32
}

type OrderShippingAddress struct {
	OrderID       int64
	RecipientName string
	PhoneNumber   string
	AddressLine   string
	Subdistrict   string
	District      string
	Province      string
	ZipCode       string
}

type CreateOrderInput struct {
	CustomerID              *int32
	Channel                 string
	PlacedByAdminUserID     *int32
	CartID                  *int32
	Subtotal                int32
	DiscountAmount          int32
	ShippingFee             int32
	TotalAmount             int32
	DiscountID              *int32
	PaymentMethodTypeID     int32
	CustomerPaymentMethodID *int32
	Items                   []CreateOrderItemInput
	ShippingAddress         *CreateOrderShippingAddressInput
	Remark                  *string
}

type CreateOrderItemInput struct {
	PartID   int32
	Quantity int32
	Price    int32
}

type CreateOrderShippingAddressInput struct {
	RecipientName string
	PhoneNumber   string
	AddressLine   string
	Subdistrict   string
	District      string
	Province      string
	ZipCode       string
}
