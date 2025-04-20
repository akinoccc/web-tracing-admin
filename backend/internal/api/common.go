package api

// 错误响应
type ErrorResponse struct {
	Message string `json:"message"`
}

// 成功响应
type SuccessResponse struct {
	Message string `json:"message"`
}
