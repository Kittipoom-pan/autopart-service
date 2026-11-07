package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbOrderToOrderRes(dbOrder db.Order) *OrderRes {
	return &OrderRes{
		OrderID:         dbOrder.OrderID,
		OrderNumber:     dbOrder.OrderNumber,
		CustomerID:      utils.NullInt32ToInt32Ptr(dbOrder.CustomerID),
		Channel:         dbOrder.Channel,
		PlacedByAdminID: utils.NullInt32ToInt32Ptr(dbOrder.PlacedByAdminID),
		CartID:          utils.NullInt32ToInt32Ptr(dbOrder.CartID),
		Subtotal:        dbOrder.Subtotal,
		DiscountAmount:  dbOrder.DiscountAmount,
		ShippingFee:     dbOrder.ShippingFee,
		TotalAmount:     dbOrder.TotalAmount,
		DiscountID:      utils.NullInt32ToInt32Ptr(dbOrder.DiscountID),
		Status:          dbOrder.Status,
		Remark:          utils.NullStringToStringPtr(dbOrder.Remark),
		CreatedAt:       dbOrder.CreatedAt,
		UpdatedAt:       utils.NullTimeToTimePtr(dbOrder.UpdatedAt),
	}
}

func MapDbOrderByOrderNumberToOrderRes(dbOrder db.Order) *OrderRes {
	return MapDbOrderToOrderRes(dbOrder)
}

func MapDbOrdersToOrderRes(dbOrder db.Order) *OrderRes {
	return MapDbOrderToOrderRes(dbOrder)
}

func MapCreateOrderInputToParams(in *CreateOrderInput, orderNumber string, orderStatus db.OrderStatus) db.CreateOrderParams {
	return db.CreateOrderParams{
		OrderNumber:     orderNumber,
		CustomerID:      utils.Int32PtrToNullInt32(in.CustomerID),
		Channel:         db.OrderChannel(in.Channel),
		PlacedByAdminID: utils.Int32PtrToNullInt32(in.PlacedByAdminID),
		CartID:          utils.Int32PtrToNullInt32(in.CartID),
		Subtotal:        in.Subtotal,
		DiscountAmount:  in.DiscountAmount,
		ShippingFee:     in.ShippingFee,
		TotalAmount:     in.TotalAmount,
		DiscountID:      utils.Int32PtrToNullInt32(in.DiscountID),
		Status:          orderStatus,
		Remark:          utils.StringPtrToNullString(in.Remark),
	}
}

func MapCreateOrderItemInputToParams(item *CreateOrderItemInput) db.CreateOrderItemParams {
	return db.CreateOrderItemParams{
		OrderID:  0, // จะถูก set ใน repository
		PartID:   item.PartID,
		Quantity: item.Quantity,
		Price:    item.Price,
	}
}

func MapCreateOrderShippingAddressInputToParams(orderID int64, address *CreateOrderShippingAddressInput) db.CreateOrderShippingAddressParams {
	return db.CreateOrderShippingAddressParams{
		OrderID:       orderID,
		RecipientName: address.RecipientName,
		PhoneNumber:   address.PhoneNumber,
		AddressLine:   address.AddressLine,
		Subdistrict:   address.Subdistrict,
		District:      address.District,
		Province:      address.Province,
		ZipCode:       address.ZipCode,
	}
}
