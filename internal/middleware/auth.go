package middleware

import (
	"context"
	"net/http"
	"strings"

	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/golang-jwt/jwt/v5"
)

const UserCtxKey = "user"

func AuthMiddleware(next http.Handler, userService services.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			apiError := handlersutils.NewAPIError(
				"MISSING_TOKEN",
				"Missing Token.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ValidateJWT(tokenString)
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

		ctx := context.WithValue(r.Context(), UserCtxKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
