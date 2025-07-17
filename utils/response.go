package utils

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type EmptyObj struct{}

func BuildResponseSuccess(status int, message string, data any) Response {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return res
}

func BuildResponseFailed(status int, message string, err string, data any) Response {
	res := Response{
		Status:  status,
		Message: message,
		Error:   err,
		Data:    data,
	}
	return res
}
