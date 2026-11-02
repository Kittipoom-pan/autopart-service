package entity

type CustomerCreateOrderReq struct {
	CartID                  int32   `json:"cart_id" validate:"required,gt=0"`
	AddressID               int32   `json:"address_id" validate:"required,gt=0"`
	PaymentMethodTypeID     int32   `json:"payment_method_type_id" validate:"required,gt=0"`
	CustomerPaymentMethodID *int32  `json:"customer_payment_method_id"`
	DiscountID              *int32  `json:"discount_id"`
	Remark                  *string `json:"remark"`
}

type AdminCreateOrderItemReq struct {
	PartID   int32 `json:"part_id" validate:"required,gt=0"`
	Quantity int32 `json:"quantity" validate:"required,gt=0"`
}

type AdminShippingAddressReq struct {
	RecipientName string `json:"recipient_name" validate:"required"`
	PhoneNumber   string `json:"phone_number" validate:"required"`
	AddressLine   string `json:"address_line" validate:"required"`
	SubdistrictID int32  `json:"subdistrict_id" validate:"required,gt=0"`
}

type AdminCreateOrderReq struct {
	CustomerID              *int32                    `json:"customer_id"`
	Channel                 string                    `json:"channel" validate:"required"`
	AddressID               *int32                    `json:"address_id"`
	PaymentMethodTypeID     int32                     `json:"payment_method_type_id" validate:"required,gt=0"`
	CustomerPaymentMethodID *int32                    `json:"customer_payment_method_id"`
	DiscountID              *int32                    `json:"discount_id"`
	Items                   []AdminCreateOrderItemReq `json:"items" validate:"required,min=1,dive"`
	ShippingAddress         *AdminShippingAddressReq  `json:"shipping_address"`
	Remark                  *string                   `json:"remark"`
}
