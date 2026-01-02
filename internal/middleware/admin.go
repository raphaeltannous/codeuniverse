package middleware

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, ok := ctx.Value(UserAuthCtxKey).(*models.User)
		if !ok {
			handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
			return
		}

		if user.Role != "admin" {
			apiError := handlersutils.NewAPIError(
				"USER_NOT_ADMIN",
				"Not an admin",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
