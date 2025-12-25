package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func staticRouter(
	staticHandler *handlers.StaticHandler,

	authMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Get("/avatars/{filename}", staticHandler.GetAvatar)
	r.Get("/courses/thumbnails/{filename}", staticHandler.GetCourseThumbnail)

	return r
}
