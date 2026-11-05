package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	addressRepository "github.com/Kittipoom-pan/autopart-service/internal/module/address/repository"
	cartRepository "github.com/Kittipoom-pan/autopart-service/internal/module/cart/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/order/repository"
	partEntity "github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
	partRepository "github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	paymentEntity "github.com/Kittipoom-pan/autopart-service/internal/module/payment/entity"
	paymentUsecase "github.com/Kittipoom-pan/autopart-service/internal/module/payment/usecase"
	pricingEntity "github.com/Kittipoom-pan/autopart-service/internal/module/pricing/entity"
	pricingUsecase "github.com/Kittipoom-pan/autopart-service/internal/module/pricing/usecase"
	shippingEntity "github.com/Kittipoom-pan/autopart-service/internal/module/shipping/entity"
	stockEntity "github.com/Kittipoom-pan/autopart-service/internal/module/stock/entity"
	stockUsecase "github.com/Kittipoom-pan/autopart-service/internal/module/stock/usecase"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type orderUsecase struct {
	repo           repository.OrderRepository
	cartRepo       cartRepository.CartRepository
	partRepo       partRepository.PartRepository
	addressRepo    addressRepository.AddressRepository
	stockUsecase   stockUsecase.StockUsecase
	paymentUsecase paymentUsecase.PaymentUsecase
	txManager      database.TxManager
	pricingUsecase pricingUsecase.PricingUsecase
	logger         zerolog.Logger
	loc            *time.Location
}

func NewOrderUsecase(
	repo repository.OrderRepository,
	cartRepo cartRepository.CartRepository,
	partRepo partRepository.PartRepository,
	addressRepo addressRepository.AddressRepository,
	stockUsecase stockUsecase.StockUsecase,
	paymentUsecase paymentUsecase.PaymentUsecase,
	txManager database.TxManager,
	pricingUsecase pricingUsecase.PricingUsecase,
	logger zerolog.Logger,
) OrderUsecase {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.FixedZone("ICT", 7*60*60)
	}
	return &orderUsecase{
		repo:           repo,
		cartRepo:       cartRepo,
		partRepo:       partRepo,
		addressRepo:    addressRepo,
		stockUsecase:   stockUsecase,
		paymentUsecase: paymentUsecase,
		txManager:      txManager,
		pricingUsecase: pricingUsecase,
		logger:         logger,
		loc:            loc,
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

func (u *orderUsecase) generateOrderNo(ctx context.Context, now time.Time) (string, error) {
	th := now.In(u.loc)
	seqDate := time.Date(th.Year(), th.Month(), th.Day(), 0, 0, 0, 0, time.UTC)

	seq, err := u.repo.NextOrderSequence(ctx, seqDate)
	if err != nil {
		return "", err
	}
	const maxDailySeq = 9999

	if seq > maxDailySeq {
		return "", fmt.Errorf("order sequence exceeded %d for %s", maxDailySeq, th.Format(time.DateOnly))
	}
	return fmt.Sprintf("OD%s%04d", th.Format("060102"), seq), nil
}

func (u *orderUsecase) CreateOrderByCustomer(ctx context.Context, req *entity.CustomerCreateOrderReq, customerID int32) (string, error) {
	// cart ต้องเป็นของ customer คนนี้ (กัน IDOR)
	cart, err := u.cartRepo.GetActiveCartByCustomerID(ctx, customerID)
	if err != nil {
		return "", helper.MapDBErrorToAPIError(err, "Cart")
	}
	if cart.CartID != req.CartID {
		return "", customerror.NewAPIError(common.StatusNotFound, "cart not found")
	}
	if len(cart.Items) == 0 {
		return "", customerror.InvalidRequestData(map[string]string{"cart": "cart is empty"})
	}

	// address ต้องเป็นของ customer คนนี้
	addr, err := u.addressRepo.GetAddressByIDAndCustomer(ctx, req.AddressID, customerID)
	if err != nil {
		return "", helper.MapDBErrorToAPIError(err, "Address")
	}

	// ใช้แค่ PartID + Quantity จาก cart, ราคาจะถูกเติมใน createOrder
	items := make([]entity.CreateOrderItemInput, len(cart.Items))
	for i, ci := range cart.Items {
		items[i] = entity.CreateOrderItemInput{PartID: ci.PartID, Quantity: ci.Quantity}
	}

	return u.createOrder(ctx, &entity.CreateOrderInput{
		CustomerID:              &customerID,
		Channel:                 string(entity.OrderChannelWeb),
		CartID:                  &req.CartID,
		PaymentMethodTypeID:     req.PaymentMethodTypeID,
		CustomerPaymentMethodID: req.CustomerPaymentMethodID,
		DiscountID:              req.DiscountID,
		Items:                   items,
		ShippingAddress: &entity.CreateOrderShippingAddressInput{
			RecipientName: addr.RecipientName,
			PhoneNumber:   addr.PhoneNumber,
			AddressLine:   addr.AddressLine,
			Subdistrict:   addr.Subdistrict,
			District:      addr.District,
			Province:      addr.Province,
			ZipCode:       addr.ZipCode,
		},
		Remark: req.Remark,
	}, &shippingEntity.ShippingDestination{ProvinceID: addr.ProvinceID})
}

func (u *orderUsecase) CreateOrderByAdmin(ctx context.Context, req *entity.AdminCreateOrderReq, adminID int32) (string, error) {
	if !entity.OrderChannel(req.Channel).IsValidOrderChannel() {
		return "", customerror.InvalidRequestData(map[string]string{"channel": "invalid channel"})
	}
	if req.ShippingAddress == nil {
		return "", customerror.InvalidRequestData(map[string]string{"shipping_address": "required"})
	}

	items := make([]entity.CreateOrderItemInput, len(req.Items))
	for i, it := range req.Items {
		items[i] = entity.CreateOrderItemInput{PartID: it.PartID, Quantity: it.Quantity}
	}

	// หน้าบ้านส่งแค่ subdistrict_id ส่วนชื่อ/zip/province ระบบดึงจาก master table เอง
	detail, err := u.addressRepo.GetSubdistrictDetailByID(ctx, req.ShippingAddress.SubdistrictID)
	if err != nil {
		return "", helper.MapDBErrorToAPIError(err, "Subdistrict")
	}

	return u.createOrder(ctx, &entity.CreateOrderInput{
		CustomerID:              req.CustomerID,
		Channel:                 req.Channel,
		PlacedByAdminUserID:     &adminID,
		PaymentMethodTypeID:     req.PaymentMethodTypeID,
		CustomerPaymentMethodID: req.CustomerPaymentMethodID,
		DiscountID:              req.DiscountID,
		Items:                   items,
		ShippingAddress: &entity.CreateOrderShippingAddressInput{
			RecipientName: req.ShippingAddress.RecipientName,
			PhoneNumber:   req.ShippingAddress.PhoneNumber,
			AddressLine:   req.ShippingAddress.AddressLine,
			Subdistrict:   detail.Subdistrict,
			District:      detail.District,
			Province:      detail.Province,
			ZipCode:       detail.ZipCode,
		},
		Remark: req.Remark,
	}, &shippingEntity.ShippingDestination{ProvinceID: &detail.ProvinceID})
}

type orderPricing struct {
	Subtotal       int32
	DiscountAmount int32
	ShippingFee    int32
	TotalAmount    int32
}

// รวม part ซ้ำ + validate quantity + เรียงตาม PartID
// การเรียงลำดับทำให้ทุก order ล็อกแถว part ในลำดับเดียวกันเสมอ กัน deadlock ตอนตัดสต็อก
func normalizeItems(items []entity.CreateOrderItemInput) ([]entity.CreateOrderItemInput, error) {
	if len(items) == 0 {
		return nil, customerror.InvalidRequestData(map[string]string{"items": "required"})
	}

	idx := make(map[int32]int, len(items))
	out := make([]entity.CreateOrderItemInput, 0, len(items))

	for _, it := range items {
		if it.Quantity <= 0 {
			return nil, customerror.InvalidRequestData(map[string]string{"quantity": "must be greater than 0"})
		}
		if i, ok := idx[it.PartID]; ok {
			out[i].Quantity += it.Quantity
			continue
		}
		idx[it.PartID] = len(out)
		out = append(out, entity.CreateOrderItemInput{PartID: it.PartID, Quantity: it.Quantity})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].PartID < out[j].PartID })
	return out, nil
}

func (u *orderUsecase) fillItemPrices(ctx context.Context, items []entity.CreateOrderItemInput) (map[int32]*partEntity.PartRes, error) {
	ids := make([]int32, len(items))
	for i, it := range items {
		ids[i] = it.PartID
	}

	parts, err := u.partRepo.GetPartsByIDs(ctx, ids)
	if err != nil {
		return nil, helper.MapDBErrorToAPIError(err, "Part")
	}

	for i := range items {
		p, ok := parts[items[i].PartID]
		if !ok {
			return nil, customerror.InvalidRequestData(map[string]string{
				"part_id": fmt.Sprintf("part %d not found", items[i].PartID),
			})
		}
		items[i].Price = p.Price
	}
	return parts, nil
}

func (u *orderUsecase) calculateOrderPricing(
	ctx context.Context,
	items []entity.CreateOrderItemInput,
	parts map[int32]*partEntity.PartRes,
	discountID *int32,
	shippingDest *shippingEntity.ShippingDestination,
) (orderPricing, error) {
	pricingItems := make([]pricingEntity.PricingItemInput, len(items))
	for i, item := range items {
		p := parts[item.PartID]
		pricingItems[i] = pricingEntity.PricingItemInput{
			PartID:   item.PartID,
			Name:     p.Name,
			SKU:      p.SKU,
			Quantity: item.Quantity,
			Price:    item.Price,
			WeightG:  p.WeightG,
		}
	}

	resp, err := u.pricingUsecase.CalculateCheckoutPricing(ctx, pricingItems, discountID, shippingDest)
	if err != nil {
		return orderPricing{}, err
	}

	// ยังไม่รองรับส่วนลด: pricing คำนวณยอดสุทธิจาก subtotal + shipping เท่านั้น
	pricing := orderPricing{
		Subtotal:       resp.Subtotal,
		DiscountAmount: 0,
		ShippingFee:    resp.ShippingFee,
		TotalAmount:    resp.TotalAmount,
	}

	// กันยอดไม่สมดุล
	if pricing.TotalAmount != pricing.Subtotal-pricing.DiscountAmount+pricing.ShippingFee {
		return orderPricing{}, fmt.Errorf(
			"pricing mismatch: subtotal=%d discount=%d shipping=%d total=%d",
			pricing.Subtotal, pricing.DiscountAmount, pricing.ShippingFee, pricing.TotalAmount,
		)
	}
	return pricing, nil
}

func (u *orderUsecase) createOrder(
	ctx context.Context,
	order *entity.CreateOrderInput,
	shippingDest *shippingEntity.ShippingDestination,
) (string, error) {
	if shippingDest == nil {
		return "", customerror.InvalidRequestData(map[string]string{"shipping_address": "required"})
	}
	if order.PaymentMethodTypeID <= 0 {
		return "", customerror.InvalidRequestData(map[string]string{"payment_method_type_id": "required"})
	}

	items, err := normalizeItems(order.Items)
	if err != nil {
		return "", err
	}
	order.Items = items

	var (
		orderNo string
		orderID int64
	)

	err = u.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		// ปิด cart เป็นขั้นแรก: UPDATE ... WHERE is_checkout = 0 จะล็อกแถว cart
		// request ที่ซ้อนกันจะได้ 0 แถวและ fail ทันที ก่อนไปแตะสต็อกหรือสร้าง order
		if order.CartID != nil {
			if err := u.cartRepo.MarkCartCheckedOut(ctx, *order.CartID); err != nil {
				if errors.Is(err, cartRepository.ErrCartAlreadyCheckedOut) {
					return customerror.NewAPIError(common.StatusConflict, "cart already checked out")
				}
				return helper.MapDBErrorToAPIError(err, "Cart")
			}
		}

		parts, err := u.fillItemPrices(ctx, order.Items)
		if err != nil {
			return err
		}

		pricing, err := u.calculateOrderPricing(ctx, order.Items, parts, order.DiscountID, shippingDest)
		if err != nil {
			return err
		}
		order.Subtotal = pricing.Subtotal
		order.DiscountAmount = pricing.DiscountAmount
		order.ShippingFee = pricing.ShippingFee
		order.TotalAmount = pricing.TotalAmount

		no, err := u.generateOrderNo(ctx, utils.NowUTC())
		if err != nil {
			return err
		}

		id, err := u.repo.CreateOrder(ctx, order, no)
		if err != nil {
			return helper.MapDBErrorToAPIError(err, "Order")
		}

		for i := range order.Items {
			if err := u.repo.CreateOrderItem(ctx, id, &order.Items[i]); err != nil {
				return helper.MapDBErrorToAPIError(err, "OrderItem")
			}
		}

		if order.ShippingAddress != nil {
			if err := u.repo.CreateOrderShippingAddress(ctx, id, order.ShippingAddress); err != nil {
				return helper.MapDBErrorToAPIError(err, "OrderShippingAddress")
			}
		}

		// ตัดสต็อก stock_movement (reason=sale) ต่อ item
		// items ถูกเรียงตาม PartID แล้ว จึงล็อกแถว part ในลำดับเดียวกันทุก order
		stockItems := make([]stockEntity.StockDeductItem, len(order.Items))
		for i, item := range order.Items {
			stockItems[i] = stockEntity.StockDeductItem{
				PartID:   item.PartID,
				Quantity: item.Quantity,
				Price:    item.Price,
			}
		}

		if err := u.stockUsecase.DeductStockForOrder(ctx, id, stockItems, order.PlacedByAdminUserID); err != nil {
			return err
		}

		if err := u.paymentUsecase.CreatePendingPayment(ctx, &paymentEntity.CreatePaymentInput{
			OrderID:                 id,
			Amount:                  order.TotalAmount,
			CustomerID:              order.CustomerID,
			PaymentMethodTypeID:     order.PaymentMethodTypeID,
			CustomerPaymentMethodID: order.CustomerPaymentMethodID,
			CreatedBy:               order.PlacedByAdminUserID,
		}); err != nil {
			return err
		}

		orderNo = no
		orderID = id
		return nil
	})
	if err != nil {
		return "", err
	}

	u.logger.Info().
		Str("order_number", orderNo).
		Int64("order_id", orderID).
		Int32("subtotal", order.Subtotal).
		Int32("discount_amount", order.DiscountAmount).
		Int32("shipping_fee", order.ShippingFee).
		Int32("total_amount", order.TotalAmount).
		Msg("order created successfully")

	return orderNo, nil
}
