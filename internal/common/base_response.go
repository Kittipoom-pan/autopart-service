package common

type BaseResponse struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type BaseErrorResponse struct {
	Message string `json:"message"`
}
