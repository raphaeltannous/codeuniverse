package handlers

import (
	"encoding/json"
	"net/http"
)

func decodeJSONRequest(w http.ResponseWriter, r *http.Request, requestBody any) bool {
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		apiError := NewAPIError(
			"INVALID_REQUEST_BODY",
			"Invalid request body.",
		)

		writeResponseJSON(w, apiError, http.StatusBadRequest)

		return false
	}

	return true
}
