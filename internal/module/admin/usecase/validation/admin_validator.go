package validation

import (
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
	adminerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

func ValidateAdminRequest(req *entity.AdminReq, isUpdate bool) error {
	if req.Username == "" {
		return adminerror.NewAPIError(common.StatusBadRequest, "username is required")
	}

	// UpdateAdmin SQL always overwrites password, so it is required on create and update.
	if req.Password == "" {
		return adminerror.NewAPIError(common.StatusBadRequest, "password is required")
	}

	if req.Role == "" {
		return adminerror.NewAPIError(common.StatusBadRequest, "role is required")
	}

	if req.Role != "super_admin" && req.Role != "staff" {
		return adminerror.NewAPIError(common.StatusBadRequest, "role must be super_admin or staff")
	}

	return nil
}
