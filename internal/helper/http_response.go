package helper

import (
	"github.com/Kittipoom-pan/autopart-service/internal/common"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
)

func RespondError(c *fiber.Ctx, err error) error {
	if apiErr, ok := err.(customerror.APIError); ok {
		return c.Status(apiErr.Code).JSON(common.BaseErrorResponse{
			Message: apiErr.Message,
			Errors:  apiErr.Errors,
		})
	}

	if notFoundErr, ok := err.(*customerror.NotFoundError); ok {
		return c.Status(common.StatusNotFound).JSON(common.BaseErrorResponse{
			Message: notFoundErr.Error(),
		})
	}

	if unauthorizedErr, ok := err.(*customerror.UnauthorizedError); ok {
		return c.Status(common.StatusUnauthorized).JSON(common.BaseErrorResponse{
			Message: unauthorizedErr.Error(),
		})
	}

	if forbiddenErr, ok := err.(*customerror.ForbiddenError); ok {
		return c.Status(fiber.StatusForbidden).JSON(common.BaseErrorResponse{
			Message: forbiddenErr.Error(),
		})
	}

	return c.Status(common.StatusError).JSON(common.BaseErrorResponse{
		Message: "Internal server error",
	})
}

func RespondSuccess(c *fiber.Ctx, status int, data interface{}, message string) error {
	return c.Status(status).JSON(common.BaseResponse{
		Message: message,
		Result:  data,
	})
}
