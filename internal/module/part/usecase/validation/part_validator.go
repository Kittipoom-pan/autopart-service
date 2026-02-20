package validation

import (
	"errors"

	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entitie"
)

func ValidatePartRequest(req *entitie.PartReq, isUpdate bool) error {
	if req.Name == "" {
		return errors.New("part name cannot be empty")
	}

	if req.SKU == "" {
		return errors.New("SKU is required")
	}

	if req.Price != nil && *req.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	if req.PartBrandID == 0 {
		return errors.New("part brand ID is required")
	}

	return nil
}
