package entity

type CreatePaymentInput struct {
	OrderID                 int64
	Amount                  int32
	CustomerID              *int32
	PaymentMethodTypeID     int32
	CustomerPaymentMethodID *int32
	CreatedBy               *int32
}
