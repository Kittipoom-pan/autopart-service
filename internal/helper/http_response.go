package helper

import (
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
)

func RespondError(c *fiber.Ctx, err error) error {
	// custom APIError
	if apiErr, ok := err.(customerror.APIError); ok {
		return c.Status(apiErr.Code).JSON(common.BaseErrorResponse{
			Message: apiErr.Message,
		})
	}

	// NotFoundError
	if notFoundErr, ok := err.(*customerror.NotFoundError); ok {
		return c.Status(common.StatusNotFound).JSON(common.BaseErrorResponse{
			Message: notFoundErr.Error(),
		})
	}

	// default: Internal Server Error
	return c.Status(common.StatusError).JSON(common.BaseErrorResponse{
		Message: err.Error(),
	})
}

func RespondSuccess(c *fiber.Ctx, status int, data interface{}, message string) error {
	return c.Status(status).JSON(common.BaseResponse{
		Message: message,
		Result:  data,
	})
}
