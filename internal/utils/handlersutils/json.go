package handlersutils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func DecodeJSONRequest(w http.ResponseWriter, r *http.Request, requestBody any) bool {
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		apiError := NewAPIError(
			"INVALID_REQUEST_BODY",
			"Invalid request body.",
		)

		WriteResponseJSON(w, apiError, http.StatusBadRequest)

		return false
	}

	return true
}

func WriteResponseJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		slog.Error("writeJSON encode error", "err", err)

		http.Error(w, NewInternalServerAPIError().Message, http.StatusInternalServerError)
	}
}
