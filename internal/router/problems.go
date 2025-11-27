package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func problemsRouter(authMiddleware func(next http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	// r.Get("/{problemSlug}")

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
	})

	return r
}
