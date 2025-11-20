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
	"git.riyt.dev/codeuniverse/internal/utils"
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

	decodeJSONRequest(w, r, &requestBody)

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

	// TODO should send token
	response := map[string]string{
		"username": user.Username,
	}

	writeResponseJSON(w, response, http.StatusAccepted)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decodeJSONRequest(w, r, &requestBody)

	ctx := r.Context()

	user, err := h.userService.GetByUsername(
		ctx,
		requestBody.Username,
	)

	if err != nil {
		apiError := NewAPIError(
			"INTERNAL_SERVER_ERROR",
			"Internal Server Error.",
		)

		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			apiError.Code = "INVALID_CREDENTIALS"
			apiError.Message = "Invalid Credentials."

			writeResponseJSON(w, apiError, http.StatusUnauthorized)
		default:
			writeResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	if !utils.CheckPassword(user.PasswordHash, requestBody.Password) {
		apiError := NewAPIError(
			"INVALID_CREDENTIALS",
			"Invalid Credentials.",
		)

		writeResponseJSON(w, apiError, http.StatusUnauthorized)
		return
	}

	mfaCode, token, err := h.userService.CreateMfaCodeAndToken(ctx, user)
	if err != nil {
		writeResponseJSON(w, NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err = h.userService.SendMfaCodeVerificationEmail(
		ctx,
		user,
		mfaCode,
	)
	if err != nil {
		slog.Error("login handler error: send mfa code email", "err", err)
		writeResponseJSON(w, NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"username": user.Username,
		"mfaToken": token,
	}

	writeResponseJSON(w, response, http.StatusAccepted) // TODO what should the return status be?
}

func (h *UserHandler) MfaVerification(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Token string `json:"token"`
		Code  string `json:"code"`
	}

	decodeJSONRequest(w, r, &requestBody)

	ctx := r.Context()

	mfaCode, err := h.userService.VerifyMfaCode(
		ctx,
		requestBody.Token,
		requestBody.Code,
	)

	if err != nil {
		apiError := NewInternalServerAPIError()
		switch {
		case errors.Is(err, services.ErrTimeIsExpired):
			apiError.Code = "MFA_CODE_EXPIRED"
			apiError.Message = "Time is expired."

			writeResponseJSON(w, apiError, http.StatusUnauthorized)
		case errors.Is(err, services.ErrInvalidMfaCode):
			apiError.Code = "MFA_CODE_INVALID"
			apiError.Message = "Code is expired."

			writeResponseJSON(w, apiError, http.StatusUnauthorized)
		default:
			writeResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	user, err := h.userService.GetById(
		ctx,
		mfaCode.UserId.String(),
	)
	if err != nil {
		writeResponseJSON(w, NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	jwtToken, err := utils.CreateJWT(user)
	if err != nil {
		slog.Error("error", "err", err)
		writeResponseJSON(w, NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"jwtToken": jwtToken,
	}

	writeResponseJSON(w, response, http.StatusAccepted) // TODO: what should be the reponse status?
}

func (h *UserHandler) GetUserInfoById(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Id *string `json:"id"`
	}

	decodeJSONRequest(w, r, &requestBody)

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
