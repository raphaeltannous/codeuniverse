package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func adminRouter(
	userHandler *handlers.UserHandler,

	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(authMiddleware)
	r.Use(middleware.AdminOnly)

	r.Group(func(r chi.Router) {
		r.Use(middleware.OffsetMiddleware)
		r.Use(middleware.LimitMiddleware)

		r.Get("/users", userHandler.GetAllUsers)
	})

	return r
}
