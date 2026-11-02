package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/payment/entity"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, in *entity.CreatePaymentInput) error
}
