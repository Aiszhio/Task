package transport

// ErrorResponse описывает формат ошибки в API
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

// SuccessResponse обёртка для всех успешных ответов с полем "message"
type SuccessResponse struct {
	Message interface{} `json:"message"`
}
