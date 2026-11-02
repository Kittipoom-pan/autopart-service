package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout/entity"
)

type CheckoutUsecase interface {
	PreviewCheckout(ctx context.Context, customerID int32, req *entity.CheckoutPreviewReq) (*entity.CheckoutPreviewRes, error)
}
