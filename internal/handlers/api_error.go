package handlers

type APIError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

func NewAPIError(code, message string) APIError {
	return APIError{
		Code:    code,
		Message: message,
	}
}

func NewInternalServerAPIError() APIError {
	return NewAPIError(
		"INTERNAL_SERVER_ERROR",
		"Internal server error. Please contact support.",
	)
}
