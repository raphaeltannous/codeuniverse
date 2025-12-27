package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func profileRouter(
	userHandler *handlers.UserHandler,

	authMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/", userHandler.GetUserProfile)
		r.Get("/me", userHandler.GetAuthenticatedProfile)
	})

	return r
}
