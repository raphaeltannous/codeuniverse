package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func submissionsRouter(authMiddleware func(next http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(authMiddleware)
	// r.Get("/", ) TODO: get all submissions
	// r.Get("/{submissionsId}") TODO: get a specific submissions

	return r
}
