package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.riyt.dev/codeuniverse/internal/services"
	"github.com/google/uuid"
)

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{
		UserService: s,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	ctx := r.Context()

	id, err := h.UserService.CreateUser(
		ctx,
		requestBody.Username,
		requestBody.Password,
		requestBody.Email,
	)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(struct {
		Id uuid.UUID `json:"UUID"`
	}{
		Id: id,
	})
}
