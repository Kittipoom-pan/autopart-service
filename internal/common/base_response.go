package common

type BaseResponse[T any] struct {
	Message string `json:"message" example:"success"`
	Result  T      `json:"result"`
}

type BaseErrorResponse struct {
	Message string            `json:"message" example:"Validation failed"`
	Errors  map[string]string `json:"errors,omitempty" swaggertype:"object,string" example:"id:invalid part id format"`
}

type IDResponse struct {
	ID int64 `json:"id" example:"1"`
}
