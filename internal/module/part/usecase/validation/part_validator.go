package validation

import (
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entity"
	parterror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

func ValidatePartRequest(req *entity.PartReq, isUpdate bool) error {
	if req.Name == "" {
		return parterror.NewAPIError(common.StatusBadRequest, "part name cannot be empty")
	}

	if req.SKU == "" {
		return parterror.NewAPIError(common.StatusBadRequest, "SKU is required")
	}

	if req.Price != nil && *req.Price <= 0 {
		return parterror.NewAPIError(common.StatusBadRequest, "price must be greater than zero")
	}

	if req.PartBrandID == 0 {
		return parterror.NewAPIError(common.StatusBadRequest, "part brand ID is required")
	}

	if req.PartTypeID == 0 {
		return parterror.NewAPIError(common.StatusBadRequest, "part type ID is required")
	}

	return nil
}
