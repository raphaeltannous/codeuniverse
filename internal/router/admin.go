package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func adminRouter(
	userHandler *handlers.UserHandler,
	problemHandler *handlers.ProblemHandler,
	statsHandler *handlers.StatsHandler,

	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(authMiddleware)
	r.Use(middleware.AdminOnly)

	r.Route("/dashboard", func(r chi.Router) {
		r.Get("/stats", statsHandler.GetDashboardStats)
		r.Get("/activity", statsHandler.GetRecentActivity)
		r.Get("/submissions-activities", statsHandler.GetSubmissionTrendsSample)
	})

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

			r.Get("/", problemHandler.GetProblems)
		})

		r.Post("/", problemHandler.CreateProblem)

		r.Route("/{problemSlug}", func(r chi.Router) {
			// TODO problemMiddleware

			r.Get("/", problemHandler.GetProblem)
			r.Put("/", problemHandler.UpdateProblem)
			r.Delete("/", problemHandler.DeleteProblem)
		})
	})

	return r
}
