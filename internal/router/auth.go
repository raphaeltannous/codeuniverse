package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func authRouter(
	userHandler *handlers.UserHandler,

	authMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", userHandler.Signup)
	r.Post("/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/logout", userHandler.Logout)
		r.Post("/refresh", userHandler.RefreshJWTToken)
	})

	return r
}
