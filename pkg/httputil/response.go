package httputil

type Res struct {
	Message    string      `json:"message"`
	StatusCode int
	Data       interface{} `json:"data,omitempty"`
}

func NewResponse(message string, statusCode int, data interface{}) Res {
	return Res{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}
}
