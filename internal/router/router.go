package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Service(
	userHandler *handlers.UserHandler,

	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)

	r.Mount("/api", apiRouter(
		userHandler,
		authMiddleware,
	))

	return r
}

func apiRouter(
	userHandler *handlers.UserHandler,

	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Mount("/health", heathRouter())
	r.Mount("/auth", authRouter(userHandler, authMiddleware))

	// Private
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Mount("/admin", adminRouter(userHandler))
	})

	return r
}
