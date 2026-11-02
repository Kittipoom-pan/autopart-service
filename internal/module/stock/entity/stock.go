package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
)

type StockDeductItem struct {
	PartID   int32
	Quantity int32
	Price    int32
}

type StockMovementInput struct {
	PartID         int32
	QuantityChange int32
	Price          int32
	EventType      db.StockMovementEventType
	Reason         db.StockMovementReason
	OrderID        *int64
	Remark         *string
	CreatedBy      *int32
}
