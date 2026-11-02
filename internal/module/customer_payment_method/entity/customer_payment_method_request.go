package entity

type CreateCustomerPaymentMethodReq struct {
	AccountName         *string `json:"account_name"`
	PaymentMethodTypeID int32   `json:"payment_method_type_id"`
	CardToken           *string `json:"card_token"`
	IsDefault           bool    `json:"is_default"`
}

type UpdateCustomerPaymentMethodReq struct {
	AccountName         *string `json:"account_name"`
	PaymentMethodTypeID int32   `json:"payment_method_type_id"`
	CardToken           *string `json:"card_token"`
}
