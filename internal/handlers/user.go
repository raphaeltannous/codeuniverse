package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.riyt.dev/codeuniverse/internal/middleware"
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

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
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

		http.Error(w, "failed to create user", http.StatusInternalServerError)
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

func (h *UserHandler) GetUserInfoById(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Id *string `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Id == nil {
		http.Error(w, "Invalid request body: id is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	user, err := h.UserService.GetUserInfoById(ctx, *requestBody.Id)
	if err != nil {
		if strings.Contains(err.Error(), "does not exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	offset, ok := r.Context().Value("offset").(int)
	if !ok {
		offset = middleware.OffsetDefault
	}

	limit, ok := r.Context().Value("limit").(int)
	if !ok {
		limit = middleware.LimitDefault
	}

	users, err := h.UserService.GetAllUsers(ctx, offset, limit)
	if err != nil {
		http.Error(w, "failed to fetch users"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) RefreshJWTToken(w http.ResponseWriter, r *http.Request) {

}
