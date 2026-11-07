package usecase

import (
	"context"
	"time"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	cpmRepository "github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/payment/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/payment/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type paymentUsecase struct {
	repo                      repository.PaymentRepository
	customerPaymentMethodRepo cpmRepository.CustomerPaymentMethodRepository
	logger                    zerolog.Logger
}

func NewPaymentUsecase(
	repo repository.PaymentRepository,
	customerPaymentMethodRepo cpmRepository.CustomerPaymentMethodRepository,
	logger zerolog.Logger,
) PaymentUsecase {
	return &paymentUsecase{
		repo:                      repo,
		customerPaymentMethodRepo: customerPaymentMethodRepo,
		logger:                    logger,
	}
}

// validateSavedPaymentMethod ตรวจบัตรที่ลูกค้าบันทึกไว้ (ถ้ามีการอ้างถึง)
//
//	ต้องมี CustomerID ประกอบเสมอ (ไม่ข้ามการตรวจเมื่อ CustomerID เป็น nil)
//	ชนิดต้องตรงกับ PaymentMethodTypeID ที่ส่งมา
func (u *paymentUsecase) validateSavedPaymentMethod(ctx context.Context, in *entity.CreatePaymentInput) error {
	if in.CustomerPaymentMethodID == nil {
		return nil
	}
	if in.CustomerID == nil {
		return customerror.InvalidRequestData(map[string]string{
			"customer_payment_method_id": "requires customer_id",
		})
	}

	cpm, err := u.customerPaymentMethodRepo.GetCustomerPaymentMethodByID(ctx, *in.CustomerPaymentMethodID, *in.CustomerID)
	if err != nil {
		return helper.MapDBErrorToAPIError(err, "Customer payment method")
	}

	if cpm.PaymentMethodTypeID != in.PaymentMethodTypeID {
		return customerror.InvalidRequestData(map[string]string{
			"payment_method_type_id": "does not match customer payment method",
		})
	}
	return nil
}

func (u *paymentUsecase) CreatePayment(ctx context.Context, in *entity.CreatePaymentInput, isCompletePayment bool) error {
	if err := u.validateSavedPaymentMethod(ctx, in); err != nil {
		return err
	}

	status := db.PaymentStatusPending
	var paidAt *time.Time

	if isCompletePayment {
		status = db.PaymentStatusCompleted
		now := utils.NowUTC()
		paidAt = &now
	}

	in.Status = status
	in.PaidAt = paidAt
	u.logger.Info().Interface("payment_req", in).Msg("creating payment")

	if err := u.repo.CreatePayment(ctx, in); err != nil {
		return helper.MapDBErrorToAPIError(err, "Payment")
	}

	u.logger.Info().
		Int64("order_id", in.OrderID).
		Int32("amount", in.Amount).
		Int32("payment_method_type_id", in.PaymentMethodTypeID).
		Msg("pending payment created")

	return nil
}
