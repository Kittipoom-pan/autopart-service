package entity

type CheckoutPreviewReq struct {
	CartID                  int32  `json:"cart_id"`
	AddressID               int32  `json:"address_id"`
	CustomerPaymentMethodID *int32 `json:"customer_payment_method_id,omitempty"`
	DiscountID              *int32 `json:"discount_id,omitempty"`
}
