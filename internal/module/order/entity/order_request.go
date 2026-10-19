package entity

type CreateOrderInput struct {
	CustomerID              *int32
	Channel                 string
	PlacedByAdminUserID     *int32
	CustomerPaymentMethodID *int32
	CartID                  *int32
	AddressID               *int32
	DiscountID              *int32
	Amount                  int32
	Remark                  *string
}

type CustomerCreateOrderReq struct {
	CartID                  int32   `json:"cart_id"`
	AddressID               int32   `json:"address_id"`
	CustomerPaymentMethodID *int32  `json:"customer_payment_method_id"`
	DiscountID              *int32  `json:"discount_id"`
	Remark                  *string `json:"remark"`
}

type AdminCreateOrderReq struct {
	CustomerID              *int32  `json:"customer_id"`
	Channel                 string  `json:"channel"`
	CartID                  *int32  `json:"cart_id"`
	AddressID               *int32  `json:"address_id"`
	CustomerPaymentMethodID *int32  `json:"customer_payment_method_id"`
	DiscountID              *int32  `json:"discount_id"`
	Amount                  int32   `json:"amount"`
	Remark                  *string `json:"remark"`
}
