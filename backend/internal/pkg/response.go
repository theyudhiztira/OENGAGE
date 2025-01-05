package pkg

type ApiResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  bool        `json:"status"`
}

func GenerateResponse(data interface{}, message string, status ...bool) ApiResponse {
	responseStatus := true
	if len(status) > 0 {
		responseStatus = status[0]
	}

	return ApiResponse{
		Data:    data,
		Message: message,
		Status:  responseStatus,
	}
}

const (
	BadRequest          = "Bad Request"
	Unauthorized        = "Unauthorized"
	InternalServerError = "Internal Server Error"
	Forbidden           = "Forbidden"
)
