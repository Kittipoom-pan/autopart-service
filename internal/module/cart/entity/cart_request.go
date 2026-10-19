package entity

type AddCartItemReq struct {
	PartID   int32 `json:"part_id"`
	Quantity int32 `json:"quantity"`
}

type UpdateCartItemQuantityReq struct {
	Quantity int32 `json:"quantity"`
}
