package validation

import (
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/cart/entity"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

func ValidateAddCartItemRequest(req *entity.AddCartItemReq) error {
	if req.PartID <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "part_id is required and must be greater than zero")
	}
	if req.Quantity <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "quantity must be greater than zero")
	}
	return nil
}

func ValidateUpdateCartItemQuantityRequest(cartItemID int32, req *entity.UpdateCartItemQuantityReq) error {
	if cartItemID <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "item id must be greater than zero")
	}
	if req.Quantity <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "quantity must be greater than zero")
	}
	return nil
}
