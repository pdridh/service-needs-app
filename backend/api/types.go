package api

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewSuccessResponse(status int, message string, data any) *SuccessResponse {
	return &SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
