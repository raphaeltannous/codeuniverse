package handlersutils

import "net/http"

type APISuccess struct {
	Message string `json:"message"`
}

func NewAPISuccess(message string) *APISuccess {
	return &APISuccess{
		Message: message,
	}
}

func WriteSuccessMessage(w http.ResponseWriter, message string, status int) {
	WriteResponseJSON(
		w,
		NewAPISuccess(message),
		status,
	)
}
