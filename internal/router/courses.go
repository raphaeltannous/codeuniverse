package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func courseRouter(
	coursesHandler *handlers.CourseHandler,

	authMiddleware func(next http.Handler) http.Handler,
	courseMiddleware func(next http.Handler) http.Handler,
	lessonMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/", coursesHandler.GetPublicCourses)

	r.Route("/{courseSlug}", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(courseMiddleware)

		r.Get("/", coursesHandler.GetLessons)
	})

	return r
}
