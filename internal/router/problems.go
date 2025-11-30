package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

func problemsRouter(
	problemsHandler *handlers.ProblemHanlder,
	authMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Get("/", problemsHandler.GetProblems)

	r.Route("/{problemSlug}", func(r chi.Router) {
		r.Get("/", problemsHandler.GetProblem)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Route("/submit", func(r chi.Router) {
				r.Post("/", handlersutils.Unimplemented)

				r.Get("/{submissionId}/check", handlersutils.Unimplemented)
			})

			r.Route("/notes", func(r chi.Router) {
				r.Get("/", handlersutils.Unimplemented)
				r.Post("/", handlersutils.Unimplemented)
				r.Put("/", handlersutils.Unimplemented)
				r.Delete("/", handlersutils.Unimplemented)
			})
		})
	})

	return r
}
