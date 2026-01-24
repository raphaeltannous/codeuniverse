package middleware

import (
	"context"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/golang-jwt/jwt/v5"
)

const UserAuthCtxKey = "userAuth"

func PartialAuthMiddleware(next http.Handler, userService services.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("jwt")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		AuthMiddleware(next, userService).ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler, userService services.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			apiError := handlersutils.NewAPIError(
				"UNAUTHORIZED",
				"Unauthorized.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
			return
		}

		token, err := utils.ValidateJWT(cookie.Value)
		if err != nil || !token.Valid {
			apiError := handlersutils.NewAPIError(
				"INVALID_TOKEN",
				"Invalid Token.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			apiError := handlersutils.NewAPIError(
				"INVALID_Claims",
				"Invalid Claims.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
			return
		}

		user, err := userService.GetByUsername(r.Context(), username)
		if err != nil {
			apiError := handlersutils.NewAPIError(
				"INVALID_TOKEN_USER",
				"User not found.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
			return
		}

		if !user.IsActive {
			apiError := handlersutils.NewAPIError(
				"USER_IS_NOT_ACTIVE",
				"User banned or disabled.",
			)

			handlersutils.RemoveJwtCookie(w)
			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserAuthCtxKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func UserVerified(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, ok := ctx.Value(UserAuthCtxKey).(*models.User)
		if !ok {
			handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
			return
		}

		if !user.IsVerified {
			apiError := handlersutils.NewAPIError(
				"USER_NOT_VERIFIED",
				"User not verified.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
