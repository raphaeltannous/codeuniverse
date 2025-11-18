package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username        string `json:"username"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"confirm"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		apiError := NewAPIError(
			"INVALID_REQUEST_BODY",
			"Invalid request body.",
		)

		writeResponseJSON(w, apiError, http.StatusBadRequest)
		return
	}

	if requestBody.Password != requestBody.PasswordConfirm {
		apiError := NewAPIError(
			"PASSWORD_MISMATCH",
			"Passwords do not match.",
		)

		writeResponseJSON(w, apiError, http.StatusConflict)
		return
	}

	ctx := r.Context()

	user, err := h.userService.Create(
		ctx,
		requestBody.Username,
		requestBody.Password,
		requestBody.Email,
	)

	if err != nil {
		apiError := NewAPIError(
			"INTERNAL_SERVER_ERROR",
			"Internal server error. Please contact support.",
		)

		switch {
		case errors.Is(err, repository.ErrUserAlreadyExists):
			apiError.Code = "USER_ALREADY_EXISTS"
			apiError.Message = "User already exists."

			writeResponseJSON(w, apiError, http.StatusConflict)

		default:
			writeResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	response := map[string]string{
		"username": user.Username,
	}

	writeResponseJSON(w, response, http.StatusAccepted)
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

	user, err := h.userService.GetById(ctx, *requestBody.Id)
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

	users, err := h.userService.GetAllUsers(ctx, offset, limit)
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

func (h *UserHandler) PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err = h.userService.SendPasswordResetEmail(
		ctx,
		requestBody.Email,
	)
	if err != nil {
		http.Error(w, "failed to send email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	fmt.Fprint(w, "email is sent")
}

func (h *UserHandler) PasswordResetByToken(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"confirm"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Password != requestBody.PasswordConfirm {
		http.Error(w, "passwords do not match", http.StatusConflict)
		return
	}

	ctx := r.Context()

	err = h.userService.ResetPasswordByToken(
		ctx,
		requestBody.Token,
		requestBody.Password,
	)

	if err != nil {
		slog.Error("failed to reset password", "err", err)
		http.Error(w, "failed to reset password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	fmt.Fprint(w, "password is reset")

}

func (h *UserHandler) VerifyEmailByToken(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	time.Sleep(3 * time.Second)

	ctx := r.Context()

	err = h.userService.VerifyEmailByToken(
		ctx,
		requestBody.Token,
	)

	if err != nil {
		slog.Error("failed to verify email", "err", err)
		http.Error(w, "failed to verify email"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	fmt.Fprint(w, "email is verified")
}
