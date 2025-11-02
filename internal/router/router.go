package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Service(userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)

	r.Mount("/api", apiRouter(
		userHandler,
	))

	return r
}

func apiRouter(userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Mount("/health", heathRouter())
	r.Mount("/auth", authRouter(userHandler))

	// Private
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Mount("/admin", adminRouter())
	})

	return r
}
