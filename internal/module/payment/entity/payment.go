package entity

import (
	"time"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
)

type CreatePaymentInput struct {
	OrderID                 int64
	Amount                  int32
	CustomerID              *int32
	PaymentMethodTypeID     int32
	CustomerPaymentMethodID *int32
	Status                  db.PaymentStatus
	CreatedBy               *int32
	PaidAt                  *time.Time
}
