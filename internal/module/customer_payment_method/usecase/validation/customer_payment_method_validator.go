package validation

import (
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer_payment_method/entity"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

func ValidateCreateCustomerPaymentMethodRequest(req *entity.CreateCustomerPaymentMethodReq) error {
	if req.PaymentMethodTypeID <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "payment_method_type_id must be greater than zero")
	}
	return nil
}

func ValidateUpdateCustomerPaymentMethodRequest(id int32, req *entity.UpdateCustomerPaymentMethodReq) error {
	if id <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "customer_payment_method_id must be greater than zero")
	}
	if req.PaymentMethodTypeID <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "payment_method_type_id must be greater than zero")
	}
	return nil
}
