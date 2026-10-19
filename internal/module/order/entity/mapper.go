package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

// ใช้ db.Order เป็นตัวรับข้อมูลกลางสำหรับทุกฟังก์ชัน Query เนื่องจาก SQLC ใช้โมเดลเดียวกันทั้งหมดแล้ว
func MapDbOrderToOrderRes(dbOrder db.Order) *OrderRes {
	return &OrderRes{
		OrderID:                 int32(dbOrder.OrderID),
		OrderNumber:             dbOrder.OrderNumber,
		CustomerID:              utils.NullInt32ToInt32Ptr(dbOrder.CustomerID), // แปลง sql.NullInt32 เป็น *int32 ตามโครงสร้าง Response ของคุณ
		Channel:                 string(dbOrder.Channel),
		PlacedByAdminUserID:     utils.NullInt32ToInt32Ptr(dbOrder.PlacedByAdminUserID),
		CustomerPaymentMethodID: utils.NullInt32ToInt32Ptr(dbOrder.CustomerPaymentMethodID),
		CartID:                  utils.NullInt32ToInt32Ptr(dbOrder.CartID),
		Amount:                  dbOrder.Amount,
		AddressID:               utils.NullInt32ToInt32Ptr(dbOrder.AddressID),
		DiscountID:              utils.NullInt32ToInt32Ptr(dbOrder.DiscountID),
		Status:                  string(dbOrder.Status),
		Remark:                  utils.NullStringToStringPtr(dbOrder.Remark), // แปลง sql.NullString เป็น *string
		CreatedAt:               dbOrder.CreatedAt,
		UpdatedAt:               utils.NullTimeToTimePtr(dbOrder.UpdatedAt), // แปลง sql.NullTime เป็น *time.Time
	}
}

// ฟังก์ชันเดิมสำหรับดึงตาม ID
func MapDbOrderRes(dbOrder db.Order) *OrderRes {
	return MapDbOrderToOrderRes(dbOrder)
}

// ฟังก์ชันดึงตาม Order Number สามารถเรียกใช้ตัวกลางร่วมกันได้เลย
func MapDbOrderByOrderNumberToOrderRes(dbOrder db.Order) *OrderRes {
	return MapDbOrderToOrderRes(dbOrder)
}

// ฟังก์ชันดึงแบบ List สามารถเรียกใช้ตัวกลางร่วมกันได้เลย
func MapDbOrdersToOrderRes(dbOrder db.Order) *OrderRes {
	return MapDbOrderToOrderRes(dbOrder)
}

func MapOrderToCreateOrderParams(req *OrderReq, adminUserID *int32, orderNumber string) db.CreateOrderParams {
	return db.CreateOrderParams{
		OrderNumber:             orderNumber,
		CustomerID:              utils.IntNullToNullInt32(req.CustomerID),
		Channel:                 db.OrderChannel(req.Channel),
		PlacedByAdminUserID:     utils.Int32PtrToNullInt32(adminUserID),
		CustomerPaymentMethodID: utils.IntNullToNullInt32(req.CustomerPaymentMethodID),
		CartID:                  utils.IntNullToNullInt32(req.CartID),
		Amount:                  int32(req.Amount),
		AddressID:               utils.IntNullToNullInt32(req.AddressID),
		DiscountID:              utils.IntNullToNullInt32(req.DiscountID),
		Status:                  db.OrderStatusPending,
		Remark:                  utils.StringPtrToNullString(req.Remark),
		CreatedAt:               utils.NowUTC(),
	}
}

func MapCreateOrderInputToParams(in *CreateOrderInput, orderNumber string) db.CreateOrderParams {
	return db.CreateOrderParams{
		OrderNumber:             orderNumber,
		CustomerID:              utils.Int32PtrToNullInt32(in.CustomerID),
		Channel:                 db.OrderChannel(in.Channel),
		PlacedByAdminUserID:     utils.Int32PtrToNullInt32(in.PlacedByAdminUserID),
		CustomerPaymentMethodID: utils.Int32PtrToNullInt32(in.CustomerPaymentMethodID),
		CartID:                  utils.Int32PtrToNullInt32(in.CartID),
		Amount:                  in.Amount,
		AddressID:               utils.Int32PtrToNullInt32(in.AddressID),
		DiscountID:              utils.Int32PtrToNullInt32(in.DiscountID),
		Status:                  db.OrderStatusPending,
		Remark:                  utils.StringPtrToNullString(in.Remark),
		CreatedAt:               utils.NowUTC(),
	}
}
