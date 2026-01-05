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
	UserStatusFilterCtxKey       = "userStatusFilter"
	UserVerificationFilterCtxKey = "userVerifiedFilter"
	UserRoleFilterCtxKey         = "userRoleFilter"
	UserSortByFilterCtxKey       = "userSortByFilter"
	UserSortOrderFilterCtxKey    = "userSortOrderFilter"

	UserCtxKey = "user"
)

var UserStatusFilterMiddleware = makeRepositoryParamMiddleware(
	"status",
	UserStatusFilterCtxKey,
	map[string]repository.UserParam{
		"active":   repository.UserActive,
		"inactive": repository.UserInactive,
	},
)

var UserVerificationFilterMiddleware = makeRepositoryParamMiddleware(
	"verified",
	UserVerificationFilterCtxKey,
	map[string]repository.UserParam{
		"verified":   repository.UserVerified,
		"unverified": repository.UserUnverified,
	},
)

var UserRoleFilterMiddleware = makeRepositoryParamMiddleware(
	"role",
	UserRoleFilterCtxKey,
	map[string]repository.UserParam{
		"user":  repository.UserRoleUser,
		"admin": repository.UserRoleAdmin,
	},
)

var UserSortByFilterMiddleware = makeRepositoryParamMiddleware(
	"sortBy",
	UserSortByFilterCtxKey,
	map[string]repository.UserParam{
		"username":  repository.UserSortByUsername,
		"email":     repository.UserSortByEmail,
		"createdAt": repository.UserSortByCreatedAt,
	},
)

var UserSortOrderFilterMiddleware = makeRepositoryParamMiddleware(
	"sortOrder",
	UserSortOrderFilterCtxKey,
	map[string]repository.UserParam{
		"desc": repository.UserSortOrderDesc,
		"asc":  repository.UserSortOrderAsc,
	},
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
