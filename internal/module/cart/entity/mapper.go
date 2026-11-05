package entity

import (
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
)

func MapDbCartToCartRes(dbCart db.GetCartByCustomerIDRow, items []*CartItemRes) *CartRes {
	if items == nil {
		items = []*CartItemRes{}
	}
	return &CartRes{
		CartID:   dbCart.CartID,
		Subtotal: dbCart.Subtotal,
		Items:    items,
	}
}

func MapDbCartItemWithPartToCartItemRes(row db.ListCartItemsWithPartRow) *CartItemRes {
	return &CartItemRes{
		CartItemID: row.CartItemID,
		PartID:     row.PartID,
		PartName:   row.PartName,
		SKU:        row.Sku,
		Price:      row.Price,
		Quantity:   row.Quantity,
		Subtotal:   row.Quantity * row.Price,
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
