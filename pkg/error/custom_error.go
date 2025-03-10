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
