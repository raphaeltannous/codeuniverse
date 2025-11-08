package middleware

import (
	"context"
	"net/http"
	"strings"

	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

const UserCtxKey = "user"

func AuthMiddleware(next http.Handler, userService services.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		username := claims["username"].(string)
		user, err := userService.GetByUsername(r.Context(), username)
		if err != nil {
			http.Error(w, "username not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
