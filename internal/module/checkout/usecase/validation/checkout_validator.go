package validation

import (
	"github.com/Kittipoom-pan/autopart-service/internal/module/checkout/entity"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

func ValidateCheckoutPreviewRequest(req *entity.CheckoutPreviewReq) error {
	if req == nil {
		return customerror.InvalidRequestData(map[string]string{"body": "request body is required"})
	}

	errors := make(map[string]string)

	if req.CartID <= 0 {
		errors["cart_id"] = "cart_id is required and must be greater than zero"
	}

	if req.AddressID <= 0 {
		errors["address_id"] = "address_id is required and must be greater than zero"
	}

	if req.DiscountID != nil && *req.DiscountID <= 0 {
		errors["discount_id"] = "discount_id must be greater than zero"
	}

	if req.CustomerPaymentMethodID != nil && *req.CustomerPaymentMethodID <= 0 {
		errors["customer_payment_method_id"] = "customer_payment_method_id must be greater than zero"
	}

	if len(errors) > 0 {
		return customerror.InvalidRequestData(errors)
	}

	return nil
}
