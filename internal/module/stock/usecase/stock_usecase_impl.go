package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/stock/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/stock/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/rs/zerolog"
)

type stockUsecase struct {
	repo   repository.StockRepository
	logger zerolog.Logger
}

func NewStockUsecase(repo repository.StockRepository, logger zerolog.Logger) StockUsecase {
	return &stockUsecase{repo: repo, logger: logger}
}

// ตัดสต็อก + บันทึก stock_movement ต่อ item ต้องเรียกภายใน transaction ของ order เท่านั้น
func (u *stockUsecase) DeductStockForOrder(ctx context.Context, orderID int64, items []entity.StockDeductItem, createdBy *int32) error {
	for _, item := range items {
		if err := u.repo.DeductStock(ctx, item.PartID, item.Quantity, createdBy); err != nil {
			if errors.Is(err, repository.ErrInsufficientStock) {
				u.logger.Warn().
					Int64("order_id", orderID).
					Int32("part_id", item.PartID).
					Int32("quantity", item.Quantity).
					Msg("insufficient stock for order")
				return customerror.NewAPIError(common.StatusConflict, fmt.Sprintf("insufficient stock for part %d", item.PartID))
			}
			return helper.MapDBErrorToAPIError(err, "Stock")
		}

		if err := u.repo.CreateStockMovement(ctx, &entity.StockMovementInput{
			PartID:         item.PartID,
			QuantityChange: -item.Quantity,
			Price:          item.Price,
			EventType:      db.StockMovementEventTypeOut,
			Reason:         db.StockMovementReasonSale,
			OrderID:        &orderID,
			CreatedBy:      createdBy,
		}); err != nil {
			return helper.MapDBErrorToAPIError(err, "Stock movement")
		}
	}

	return nil
}
