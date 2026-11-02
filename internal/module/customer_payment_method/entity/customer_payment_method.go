package entity

import "time"

type CustomerPaymentMethod struct {
	CustomerPaymentMethodID int32
	CustomerID              int32
	AccountName             *string
	PaymentMethodTypeID     int32
	CardToken               *string
	IsDefault               bool
	CreatedAt               time.Time
	CreatedBy               *int32
	UpdatedAt               *time.Time
	UpdatedBy               *int32
}
