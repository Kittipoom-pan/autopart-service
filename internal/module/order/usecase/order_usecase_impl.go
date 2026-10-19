package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/repository"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type orderUsecase struct {
	repo   repository.OrderRepository
	logger zerolog.Logger
	loc    *time.Location
}

func NewOrderUsecase(repo repository.OrderRepository, logger zerolog.Logger) OrderUsecase {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.FixedZone("ICT", 7*60*60)
	}
	return &orderUsecase{
		repo:   repo,
		logger: logger,
		loc:    loc,
	}
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, id int) (*entity.OrderRes, error) {
	u.logger.Debug().Int("order_id", id).Msg("GetOrderByID")
	order, err := u.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Order")
	}
	return order, nil
}

func (u *orderUsecase) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*entity.OrderRes, error) {
	u.logger.Debug().Str("order_number", orderNumber).Msg("GetOrderByOrderNumber")
	order, err := u.repo.GetOrderByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Order")
	}
	return order, nil
}

func (u *orderUsecase) GetAllOrders(ctx context.Context) ([]*entity.OrderRes, error) {
	u.logger.Debug().Msg("GetAllOrders")
	orders, err := u.repo.GetAllOrders(ctx)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Order")
	}
	return orders, nil
}

func (u *orderUsecase) generateOrderNo(ctx context.Context, repo repository.OrderRepository, now time.Time) (string, error) {
	th := now.In(u.loc)
	// วันตามปฏิทินไทย แต่สร้างเป็น UTC เพื่อให้ driver (loc=UTC) ส่งวันที่ไม่เลื่อน
	seqDate := time.Date(th.Year(), th.Month(), th.Day(), 0, 0, 0, 0, time.UTC)

	seq, err := repo.NextOrderSequence(ctx, seqDate)
	if err != nil {
		return "", err
	}
	const maxDailySeq = 9999

	if seq > maxDailySeq {
		return "", fmt.Errorf("order sequence exceeded %d for %s", maxDailySeq, th.Format(time.DateOnly))
	}
	return fmt.Sprintf("OD%s%04d", th.Format("060102"), seq), nil
}

func (u *orderUsecase) createOrder(ctx context.Context, in *entity.CreateOrderInput) (string, error) {
	var orderNo string

	err := u.repo.WithTx(ctx, func(tx repository.OrderRepository) error {
		no, err := u.generateOrderNo(ctx, tx, utils.NowUTC())
		if err != nil {
			return err
		}

		orderID, err := tx.CreateOrder(ctx, in, no)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Order")
		}

		orderNo = no
		u.logger.Info().Str("order_number", no).Int64("order_id", orderID).Msg("order created successfully")
		return nil
	})
	if err != nil {
		return "", err
	}
	return orderNo, nil
}

func (u *orderUsecase) CreateOrderByCustomer(ctx context.Context, req *entity.CustomerCreateOrderReq, customerID int32) (string, error) {
	// Cart เฉพาะของ customerID นี้
	// cart, err := u.cartRepo.GetCartByIDAndCustomer(ctx, req.CartID, customerID)
	// if err != nil {
	// 	return "", err
	// }
	// //เช็ค address / payment method เป็นของ customerID
	// if err := u.addressRepo.EnsureOwnedByCustomer(ctx, req.AddressID, customerID); err != nil {
	// 	return "", err
	// }
	// if req.CustomerPaymentMethodID != nil {
	// 	if err := u.paymentRepo.EnsureOwnedByCustomer(ctx, *req.CustomerPaymentMethodID, customerID); err != nil {
	// 		return "", err
	// 	}
	// }
	// // คำนวณราคาฝั่ง server (รวมส่วนลดถ้ามี และ validate เงื่อนไขส่วนลด)
	// amount, err := u.pricing.CalculateAmount(ctx, cart, req.DiscountID)
	// if err != nil {
	// 	return "", err
	// }

	return u.createOrder(ctx, &entity.CreateOrderInput{
		CustomerID:              &customerID,
		Channel:                 string(entity.OrderChannelWeb),
		PlacedByAdminUserID:     nil,
		CustomerPaymentMethodID: req.CustomerPaymentMethodID,
		CartID:                  &req.CartID,
		AddressID:               &req.AddressID,
		DiscountID:              req.DiscountID,
		Amount:                  0, // TODO: calculate amount
		Remark:                  req.Remark,
	})
}

func (u *orderUsecase) CreateOrderByAdmin(ctx context.Context, req *entity.AdminCreateOrderReq, adminID int32) (string, error) {
	if !entity.OrderChannel(req.Channel).IsValidOrderChannel() {
		return "", customerror.InvalidRequestData(map[string]string{"channel": "invalid channel"})
	}
	if req.Amount < 0 {
		return "", customerror.InvalidRequestData(map[string]string{"amount": "must not be negative"})
	}

	// ถ้าระบุ customer_id ให้เช็คว่า cart/address เป็นของลูกค้าคนนั้น
	// if req.CustomerID != nil {
	// 	if req.CartID != nil {
	// 		if _, err := u.cartRepo.GetCartByIDAndCustomer(ctx, *req.CartID, *req.CustomerID); err != nil {
	// 			return "", err
	// 		}
	// 	}
	// 	if req.AddressID != nil {
	// 		if err := u.addressRepo.EnsureOwnedByCustomer(ctx, *req.AddressID, *req.CustomerID); err != nil {
	// 			return "", err
	// 		}
	// 	}
	// }

	return u.createOrder(ctx, &entity.CreateOrderInput{
		CustomerID:              req.CustomerID,
		Channel:                 req.Channel,
		PlacedByAdminUserID:     &adminID,
		CustomerPaymentMethodID: req.CustomerPaymentMethodID,
		CartID:                  req.CartID,
		AddressID:               req.AddressID,
		DiscountID:              req.DiscountID,
		Amount:                  req.Amount,
		Remark:                  req.Remark,
	})
}
