package validation

import (
	"strings"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entity"
	customererror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

func ValidateCustomerRequest(req *entity.CustomerReq, isUpdate bool) error {
	if strings.TrimSpace(req.Username) == "" {
		return customererror.NewAPIError(common.StatusBadRequest, "username is required")
	}

	if strings.TrimSpace(req.Email) == "" {
		return customererror.NewAPIError(common.StatusBadRequest, "email is required")
	}

	if !isUpdate && req.Password == "" {
		return customererror.NewAPIError(common.StatusBadRequest, "password is required")
	}

	if strings.TrimSpace(req.FirstName) == "" {
		return customererror.NewAPIError(common.StatusBadRequest, "first name is required")
	}

	if strings.TrimSpace(req.LastName) == "" {
		return customererror.NewAPIError(common.StatusBadRequest, "last name is required")
	}

	return nil
}
