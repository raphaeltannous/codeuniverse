package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

func problemsRouter(
	problemHandler *handlers.ProblemHandler,

	authMiddleware func(next http.Handler) http.Handler,
	partialAuthMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.OffsetMiddleware)
		r.Use(middleware.LimitMiddleware)
		r.Use(middleware.SearchMiddleware)

		r.Use(middleware.ProblemPremiumFilterMiddleware)
		r.Use(middleware.ProblemDifficultyFilterMiddleware)
		r.Use(middleware.ProblemSortByFilterMiddleware)
		r.Use(middleware.ProblemSortOrderFilterMiddleware)

		r.Get("/", problemHandler.GetProblems)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/progress", problemHandler.GetProgress)
	})

	r.Route("/{problemSlug}", func(r chi.Router) {
		r.Use(partialAuthMiddleware)
		r.Use(problemMiddleware)

		r.Get("/", problemHandler.GetProblem)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Route("/submit", func(r chi.Router) {
				r.Route("/{languageSlug}", func(r chi.Router) {
					r.Use(middleware.ProblemLanguageMiddleware)

					r.Post("/", problemHandler.Submit)
				})

				r.Get("/{submissionId}/check", problemHandler.GetSubmission)
			})

			r.Route("/run", func(r chi.Router) {
				r.Route("/{languageSlug}", func(r chi.Router) {
					r.Use(middleware.ProblemLanguageMiddleware)

					r.Post("/", problemHandler.Run)
				})

				r.Get("/{runId}/check", problemHandler.GetRun)
			})

			r.Route("/submissions", func(r chi.Router) {
				r.Get("/", problemHandler.GetSubmissions)

				r.Get("/{submissionId}", handlersutils.Unimplemented)
			})

			r.Route("/notes", func(r chi.Router) {
				r.Get("/", problemHandler.GetNote)
				r.Post("/", problemHandler.CreateNote)
				r.Put("/", problemHandler.UpdateNote)
				r.Delete("/", handlersutils.Unimplemented)
			})
		})
	})

	return r
}
