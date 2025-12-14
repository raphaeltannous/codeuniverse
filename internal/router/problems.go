package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

func problemsRouter(
	problemsHandler *handlers.ProblemHandler,

	authMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Get("/", problemsHandler.GetProblems)

	r.Route("/{problemSlug}", func(r chi.Router) {
		r.Use(problemMiddleware)

		r.Get("/", problemsHandler.GetProblem)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Route("/submit", func(r chi.Router) {
				r.Post("/", problemsHandler.Submit)

				r.Get("/{submissionId}/check", problemsHandler.GetSubmission)
			})

			r.Route("/run", func(r chi.Router) {
				r.Post("/", problemsHandler.Run)

				r.Get("/{runId}/check", problemsHandler.GetRun)
			})

			r.Route("/submissions", func(r chi.Router) {
				r.Get("/", problemsHandler.GetSubmissions)

				r.Get("/{submissionId}", handlersutils.Unimplemented)
			})

			r.Route("/notes", func(r chi.Router) {
				r.Get("/", problemsHandler.GetNote)
				r.Post("/", problemsHandler.CreateNote)
				r.Put("/", problemsHandler.UpdateNote)
				r.Delete("/", handlersutils.Unimplemented)
			})
		})
	})

	return r
}
