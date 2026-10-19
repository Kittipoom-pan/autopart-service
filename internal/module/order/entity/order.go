package entity

import (
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
)

type Order struct {
	OrderID                 int32
	OrderNumber             string
	CustomerID              *int32 // เปลี่ยนเป็น Pointer เพราะเป็น NULL ได้ (walk_in)
	Channel                 db.OrderChannel
	PlacedByAdminUserID     *int32
	CustomerPaymentMethodID *int32
	CartID                  *int32
	Amount                  int32
	AddressID               *int32 // เปลี่ยนเป็น Pointer เพราะเป็น NULL ได้ (walk_in)
	DiscountID              *int32
	Status                  db.OrderStatus
	Remark                  *string
	CreatedAt               time.Time
	UpdatedAt               *time.Time
}

// ใช้สำหรับส่ง Response กลับหา Client
type OrderRes struct {
	OrderID                 int32
	OrderNumber             string
	CustomerID              *int32 // เปลี่ยนเป็น Pointer รองรับกรณี walk_in
	Channel                 string
	PlacedByAdminUserID     *int32
	CustomerPaymentMethodID *int32
	CartID                  *int32
	Amount                  int32
	AddressID               *int32 // เปลี่ยนเป็น Pointer รองรับกรณี walk_in
	DiscountID              *int32
	Status                  string
	Remark                  *string
	CreatedAt               time.Time
	UpdatedAt               *time.Time
}

// ใช้สำหรับรับ Request จาก Client ตอนสร้าง Order
type OrderReq struct {
	CustomerID              *int
	Channel                 string `json:"channel"` // web, mobile, call, line, etc.
	CustomerPaymentMethodID *int
	CartID                  *int
	Amount                  int
	AddressID               *int
	DiscountID              *int
	Remark                  *string
}
