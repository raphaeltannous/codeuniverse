package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func courseRouter(
	coursesHandler *handlers.CourseHandler,

	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/", coursesHandler.GetPublicCourses)

	return r
}
