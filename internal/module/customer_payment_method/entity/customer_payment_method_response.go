package entity

import "time"

type CustomerPaymentMethodRes struct {
	CustomerPaymentMethodID int32      `json:"customer_payment_method_id"`
	CustomerID              int32      `json:"customer_id"`
	AccountName             *string    `json:"account_name,omitempty"`
	PaymentMethodTypeID     int32      `json:"payment_method_type_id"`
	CardToken               *string    `json:"card_token,omitempty"`
	IsDefault               bool       `json:"is_default"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
}
