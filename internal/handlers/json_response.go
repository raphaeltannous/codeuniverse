package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func writeResponseJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		slog.Error("writeJSON encode error", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
