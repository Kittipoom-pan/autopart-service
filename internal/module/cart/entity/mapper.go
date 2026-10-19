package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func parseBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case []byte:
		return len(val) > 0 && val[0] == 1
	case int64:
		return val == 1
	case int32:
		return val == 1
	case int:
		return val == 1
	default:
		return false
	}
}

func MapDbCartToCartRes(dbCart db.Cart, items []*CartItemRes) *CartRes {
	if items == nil {
		items = []*CartItemRes{}
	}
	return &CartRes{
		CartID:     dbCart.CartID,
		CustomerID: dbCart.CustomerID,
		Amount:     dbCart.Subtotal,
		IsCheckout: parseBool(dbCart.IsCheckout),
		Items:      items,
		CreatedAt:  dbCart.CreatedAt,
		UpdatedAt:  utils.NullTimeToTimePtr(dbCart.UpdatedAt),
	}
}

func MapDbCartItemWithPartToCartItemRes(row db.ListCartItemsWithPartRow) *CartItemRes {
	return &CartItemRes{
		CartItemID: row.CartItemID,
		CartID:     row.CartID,
		PartID:     row.PartID,
		PartName:   row.PartName,
		SKU:        row.Sku,
		Price:      row.Price,
		Quantity:   row.Quantity,
		Subtotal:   row.Quantity * row.Price,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  utils.NullTimeToTimePtr(row.UpdatedAt),
	}
}

func MapDbCartItemsWithPartToCartItemResList(rows []db.ListCartItemsWithPartRow) []*CartItemRes {
	items := make([]*CartItemRes, 0, len(rows))
	for _, row := range rows {
		items = append(items, MapDbCartItemWithPartToCartItemRes(row))
	}
	return items
}

func MapDbCartItemToCartItemEntity(dbItem db.CartItem) *CartItem {
	return &CartItem{
		CartItemID: dbItem.CartItemID,
		CartID:     dbItem.CartID,
		PartID:     dbItem.PartID,
		Quantity:   dbItem.Quantity,
		CreatedAt:  dbItem.CreatedAt,
		UpdatedAt:  utils.NullTimeToTimePtr(dbItem.UpdatedAt),
	}
}
