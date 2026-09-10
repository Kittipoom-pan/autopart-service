package customerror

import (
	"net/http"
)

type APIError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
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
	return APIError{
		Code:    code,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{Message: message}
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Invalid request data",
		Errors:  errors,
	}
}

func InvalidJSON() APIError {
	return APIError{
		Code:    http.StatusBadRequest,
		Message: "invalid json request data",
	}
}
