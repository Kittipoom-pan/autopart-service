package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/payment/entity"
)

type PaymentUsecase interface {
	CreatePendingPayment(ctx context.Context, in *entity.CreatePaymentInput) error
}
