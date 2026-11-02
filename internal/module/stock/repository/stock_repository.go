package repository

import (
	"context"
	"errors"

	"github.com/Kittipoom-pan/autopart-service/internal/module/stock/entity"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type StockRepository interface {
	DeductStock(ctx context.Context, partID int32, quantity int32, updatedBy *int32) error
	CreateStockMovement(ctx context.Context, in *entity.StockMovementInput) error
}
