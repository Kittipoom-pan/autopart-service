package validation

import (
	"strings"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/address/entity"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

func ValidateCreateAddressRequest(req *entity.CreateAddressReq) error {
	if strings.TrimSpace(req.RecipientName) == "" {
		return customerror.NewAPIError(common.StatusBadRequest, "recipient_name is required")
	}
	if strings.TrimSpace(req.AddressLine) == "" {
		return customerror.NewAPIError(common.StatusBadRequest, "address_line is required")
	}
	if req.SubdistrictID <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "subdistrict_id is required")
	}
	if strings.TrimSpace(req.ZipCode) == "" {
		return customerror.NewAPIError(common.StatusBadRequest, "zip_code is required")
	}
	if strings.TrimSpace(req.PhoneNumber) == "" {
		return customerror.NewAPIError(common.StatusBadRequest, "phone_number is required")
	}
	return nil
}

func ValidateUpdateAddressRequest(addressID int32, req *entity.UpdateAddressReq) error {
	if addressID <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "address_id must be greater than zero")
	}
	if strings.TrimSpace(req.RecipientName) == "" {
		return customerror.NewAPIError(common.StatusBadRequest, "recipient_name is required")
	}
	if strings.TrimSpace(req.AddressLine) == "" {
		return customerror.NewAPIError(common.StatusBadRequest, "address_line is required")
	}
	if req.SubdistrictID <= 0 {
		return customerror.NewAPIError(common.StatusBadRequest, "subdistrict_id is required")
	}
	if strings.TrimSpace(req.ZipCode) == "" {
		return customerror.NewAPIError(common.StatusBadRequest, "zip_code is required")
	}
	if strings.TrimSpace(req.PhoneNumber) == "" {
		return customerror.NewAPIError(common.StatusBadRequest, "phone_number is required")
	}
	return nil
}
