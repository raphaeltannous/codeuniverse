package middleware

import (
	"context"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

const (
	UserCtxKey                   = "user"
	UserRoleFilterCtxKey         = "roleFilter"
	UserStatusFilterCtxKey       = "statusFilter"
	UserVerificationFilterCtxKey = "verificationFilter"
)

func UserMiddleware(next http.Handler, userService services.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		ctx := r.Context()
		user, err := userService.GetByUsername(ctx, username)
		if err != nil {
			apiError := handlersutils.NewInternalServerAPIError()
			switch err {
			case repository.ErrUserNotFound:
				apiError.Code = "USER_NOT_FOUND"
				apiError.Message = "User not found."

				handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			default:
				handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
			}

			return
		}

		ctx = context.WithValue(ctx, UserCtxKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func UserRoleFilterMiddleware(next http.Handler) http.Handler {
	allowedRoleFilters := map[string]bool{
		"all":   true,
		"user":  true,
		"admin": true,
		"":      true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleFilter := r.URL.Query().Get("role")
		if !allowedRoleFilters[roleFilter] {
			apiError := handlersutils.NewAPIError(
				"UNALLOWED_ROLE_FILTER",
				"The requested role filter is not allowed.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), UserRoleFilterCtxKey, roleFilter)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func UserStatusFilterMiddleware(next http.Handler) http.Handler {
	allowedStatusFilters := map[string]bool{
		"active":   true,
		"inactive": true,
		"":         true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusFilter := r.URL.Query().Get("status")
		if !allowedStatusFilters[statusFilter] {
			apiError := handlersutils.NewAPIError(
				"UNALLOWED_STATUS_FILTER",
				"The requested status filter is not allowed.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), UserStatusFilterCtxKey, statusFilter)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func UserVerificationFilterMiddleware(next http.Handler) http.Handler {
	allowedVerificationFilter := map[string]bool{
		"verified":   true,
		"unverified": true,
		"":           true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifiedFilter := r.URL.Query().Get("verified")
		if !allowedVerificationFilter[verifiedFilter] {
			apiError := handlersutils.NewAPIError(
				"UNALLOWED_VERIFIED_FILTER",
				"The requested verified filter is not allowed.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), UserVerificationFilterCtxKey, verifiedFilter)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
