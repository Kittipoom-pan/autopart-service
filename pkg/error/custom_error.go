package customerror

import (
	"net/http"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type NotFoundError struct {
	Resource string
}

type UnauthorizedError struct {
	Message string
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	if e.Message == "" {
		return "Permission denied"
	}
	return e.Message
}

func (e *UnauthorizedError) Error() string {
	if e.Message == "" {
		return "Unauthorized access"
	}
	return e.Message
}

func (e APIError) Error() string {
	return e.Message
}

func (e *NotFoundError) Error() string {
	return e.Resource + " not found"
}

func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{Resource: resource}
}

func NewAPIError(code int, message string) APIError {
	apiErr := APIError{
		Code:    code,
		Message: message,
	}
	return apiErr
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{Message: message}
}

func InvalidRequestData(errors map[string]string) APIError {
	apiErr := APIError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Invalid request data",
	}
	return apiErr
}

func InvalidJSON() APIError {
	apiErr := APIError{
		Code:    http.StatusBadRequest,
		Message: "invalid json request data",
	}
	return apiErr
}
