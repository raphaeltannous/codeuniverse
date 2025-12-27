package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
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

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	if requestBody.Password != requestBody.PasswordConfirm {
		apiError := handlersutils.NewAPIError(
			"PASSWORD_MISMATCH",
			"Passwords do not match.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusConflict)
		return
	}

	ctx := r.Context()

	user, err := h.userService.RegisterUser(
		ctx,
		requestBody.Username,
		requestBody.Password,
		requestBody.Email,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()

		switch {
		case errors.Is(err, repository.ErrUserAlreadyExists):
			apiError.Code = "USER_ALREADY_EXISTS"
			apiError.Message = "User already exists."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusConflict)

		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	jwtToken, err := utils.CreateJWT(user)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"jwtToken": jwtToken,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusAccepted)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()

	user, err := h.userService.GetByUsername(
		ctx,
		requestBody.Username,
	)

	if err != nil {
		apiError := handlersutils.NewAPIError(
			"INTERNAL_SERVER_ERROR",
			"Internal Server Error.",
		)

		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			apiError.Code = "INVALID_CREDENTIALS"
			apiError.Message = "Invalid credentials."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	if !utils.CheckPassword(user.PasswordHash, requestBody.Password) {
		apiError := handlersutils.NewAPIError(
			"INVALID_CREDENTIALS",
			"Invalid Credentials.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		return
	}

	mfaCode, token, err := h.userService.CreateMfaCodeAndToken(ctx, user)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err = h.userService.SendMfaCodeVerificationEmail(
		ctx,
		user,
		mfaCode,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"username": user.Username,
		"mfaToken": token,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusAccepted) // TODO what should the return status be?
}

func (h *UserHandler) MfaVerification(w http.ResponseWriter, r *http.Request) {
	var requestBody middleware.MfaRequestBody

	if val, ok := r.Context().Value("requestBody").(middleware.MfaRequestBody); ok {
		requestBody = val
	} else {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	mfaCode, err := h.userService.VerifyMfaCode(
		ctx,
		requestBody.Token,
		requestBody.Code,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, services.ErrTimeIsExpired):
			apiError.Code = "MFA_CODE_EXPIRED"
			apiError.Message = "Time is expired."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		case errors.Is(err, services.ErrInvalidMfaCode), errors.Is(err, repository.ErrMfaTokenNotFound):
			apiError.Code = "MFA_CODE_INVALID"
			apiError.Message = "Invalid code."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	user, err := h.userService.GetById(
		ctx,
		mfaCode.UserId.String(),
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	jwtToken, err := utils.CreateJWT(user)
	if err != nil {
		slog.Error("error", "err", err)
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"jwtToken": jwtToken,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusAccepted) // TODO: what should be the reponse status?
}

func (h *UserHandler) ResendMfaVerification(w http.ResponseWriter, r *http.Request) {
	var requestBody middleware.MfaRequestBody

	if val, ok := r.Context().Value("requestBody").(middleware.MfaRequestBody); ok {
		requestBody = val
	} else {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	mfaCode, err := h.userService.GetMfaCodeByToken(
		ctx,
		requestBody.Token,
	)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrMfaTokenNotFound):
			apiError := handlersutils.NewAPIError(
				"INVALID_MFA_TOKEN",
				"Invalid Mfa Token.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		default:
			handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		}

		return
	}

	user, err := h.userService.GetById(
		ctx,
		mfaCode.UserId.String(),
	)

	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	newCode, newToken, err := h.userService.CreateMfaCodeAndToken(
		ctx,
		user,
	)

	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err = h.userService.SendMfaCodeVerificationEmail(
		ctx,
		user,
		newCode,
	)

	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"newToken": newToken,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusAccepted) // TODO: correct status?
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
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, users, http.StatusFound)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Bye!",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusContinue)
}

func (h *UserHandler) RefreshJWTToken(w http.ResponseWriter, r *http.Request) {
	handlersutils.WriteResponseJSON(w, handlersutils.NewAPIError("NOT_IMPLEMENTED", "Not Implemented."), http.StatusAccepted)
}

func (h *UserHandler) JWTTokenStatus(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Token is valid.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusAccepted)
}

func (h *UserHandler) PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email string `json:"email"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()

	err := h.userService.SendPasswordResetEmail(
		ctx,
		requestBody.Email,
	)

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	reponse := map[string]string{
		"message": "Email is sent.",
	}

	handlersutils.WriteResponseJSON(w, reponse, http.StatusAccepted)
}

func (h *UserHandler) PasswordResetByToken(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"confirm"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	if requestBody.Password != requestBody.PasswordConfirm {
		apiError := handlersutils.NewAPIError(
			"PASSWORDS_DO_NOT_MATCH",
			"Passwords do not match.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusConflict)
		return
	}

	ctx := r.Context()

	err := h.userService.ResetPasswordByToken(
		ctx,
		requestBody.Token,
		requestBody.Password,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, services.ErrTimeIsExpired):
			apiError.Code = "TIME_EXPIRED"
			apiError.Message = "Time is expired."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		case errors.Is(err, repository.ErrPasswordResetNotFound):
			apiError.Code = "INVALID_TOKEN"
			apiError.Message = "Invalid Token."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}
		return
	}

	reponse := map[string]string{
		"message": "Password is changed.",
	}

	handlersutils.WriteResponseJSON(w, reponse, http.StatusAccepted)
}

func (h *UserHandler) VerifyEmailByToken(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Token string `json:"token"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()

	err := h.userService.VerifyEmailByToken(
		ctx,
		requestBody.Token,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, services.ErrTimeIsExpired):
			apiError.Code = "TIME_IS_EXPIRED"
			apiError.Message = "Time is expired."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		case errors.Is(err, repository.ErrEmailVerificationNotFound):
			apiError.Code = "INVALID_TOKEN"
			apiError.Message = "Invalid link."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		slog.Error("failed to verify email", "err", err)
		return
	}

	response := map[string]string{
		"message": "Email is verified.",
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusAccepted)
}

// GET
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	userProfile, err := h.userService.GetProfile(
		ctx,
		user,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, repository.ErrUserProfileNotFound):
			apiError.Code = "USER_PROFILE_NOT_FOUND"
			apiError.Message = "Failed to get user profile. Contact Support."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		}

		return
	}

	handlersutils.WriteResponseJSON(w, userProfile, http.StatusOK)
}

// GET
func (h *UserHandler) GetAuthenticatedProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	userProfile, err := h.userService.GetProfile(
		ctx,
		user,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, repository.ErrUserProfileNotFound):
			apiError.Code = "USER_PROFILE_NOT_FOUND"
			apiError.Message = "Failed to get user profile. Contact Support."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		}

		return
	}

	response := struct {
		Username   string `json:"username"`
		AvatarUrl  string `json:"avatarUrl"`
		IsVerified bool   `json:"isVerified"`
		IsActive   bool   `json:"isActive"`
		Role       string `json:"role"`
	}{
		Username:   user.Username,
		AvatarUrl:  *userProfile.AvatarURL,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		Role:       user.Role,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

// GET
func (h *UserHandler) GetPulbicUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")
	user, err := h.userService.GetByUsername(
		ctx,
		username,
	)
	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			apiError.Code = "USER_NOT_FOUND"
			apiError.Message = "User not found."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		}

		return
	}

	userProfile, err := h.userService.GetProfile(
		ctx,
		user,
	)

	if err != nil {
		apiError := handlersutils.NewInternalServerAPIError()
		switch {
		case errors.Is(err, repository.ErrUserProfileNotFound):
			apiError.Code = "USER_PROFILE_NOT_FOUND"
			apiError.Message = "Failed to get user profile. Contact Support."

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		default:
			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
		}

		return
	}

	handlersutils.WriteResponseJSON(w, userProfile, http.StatusOK)
}
