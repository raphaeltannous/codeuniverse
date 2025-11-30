package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func adminRouter(
	userHandler *handlers.UserHandler,
	problemsHandler *handlers.ProblemHandler,

	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	// r.Use(authMiddleware)
	// r.Use(middleware.AdminOnly)

	r.Group(func(r chi.Router) {
		r.Use(middleware.OffsetMiddleware)
		r.Use(middleware.LimitMiddleware)
		// TODO r.Use(middleware.SearchMiddleware)

		r.Get("/users", userHandler.GetAllUsers)

	})

	r.Route("/problems", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.OffsetMiddleware)
			r.Use(middleware.LimitMiddleware)
			// TODO r.Use(middleware.SearchMiddleware)

			r.Get("/", problemsHandler.GetProblems)
		})

		r.Post("/", problemsHandler.CreateProblem)

		r.Route("/{problemSlug}", func(r chi.Router) {
			// TODO problemMiddleware

			r.Get("/", problemsHandler.GetProblem)
			r.Put("/", problemsHandler.UpdateProblem)
			r.Delete("/", problemsHandler.DeleteProblem)
		})
	})

	return r
}
