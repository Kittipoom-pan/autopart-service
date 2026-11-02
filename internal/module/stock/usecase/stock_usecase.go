package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/stock/entity"
)

type StockUsecase interface {
	DeductStockForOrder(ctx context.Context, orderID int64, items []entity.StockDeductItem, createdBy *int32) error
}
