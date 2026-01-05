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
	adminHandler *handlers.AdminHandler,

	authMiddleware func(next http.Handler) http.Handler,
	courseMiddleware func(next http.Handler) http.Handler,
	lessonMiddleware func(next http.Handler) http.Handler,
	userMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(authMiddleware)
	r.Use(middleware.AdminOnly)

	r.Route("/dashboard", func(r chi.Router) {
		r.Get("/stats", statsHandler.GetDashboardStats)
		r.Get("/activity", statsHandler.GetRecentActivity)
		r.Get("/submissions-activities", statsHandler.GetSubmissionTrendsSample)

	})

	r.Route("/courses", func(r chi.Router) {
		r.Get("/", adminHandler.GetCourses)
		r.Post("/", adminHandler.CreateCourse)

		r.Route("/{courseSlug}", func(r chi.Router) {
			r.Use(courseMiddleware)

			r.Delete("/", adminHandler.DeleteCourse)
			r.Put("/", adminHandler.UpdateCourseInfo)

			r.Put("/publish", adminHandler.UpdateCoursePublishStatus)
			r.Put("/thumbnail", adminHandler.UpdateThumbnail)

			r.Route("/lessons", func(r chi.Router) {
				r.Get("/", adminHandler.GetLessons)
				r.Post("/", adminHandler.CreateLesson)

				r.Route("/{lessonId}", func(r chi.Router) {
					r.Use(lessonMiddleware)

					r.Put("/", adminHandler.UpdateLesson)
					r.Delete("/", adminHandler.DeleteLesson)

					r.Put("/video", adminHandler.UpdateLessonVideo)
				})
			})
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", adminHandler.AddUser)

		r.Group(func(r chi.Router) {
			r.Use(middleware.OffsetMiddleware)
			r.Use(middleware.LimitMiddleware)
			r.Use(middleware.SearchMiddleware)

			r.Use(middleware.UserRoleFilterMiddleware)
			r.Use(middleware.UserStatusFilterMiddleware)
			r.Use(middleware.UserVerificationFilterMiddleware)
			r.Use(middleware.UserSortByFilterMiddleware)
			r.Use(middleware.UserSortOrderFilterMiddleware)

			r.Get("/", adminHandler.GetUsers)
		})

		r.Route("/{username}", func(r chi.Router) {
			r.Use(userMiddleware)

			r.Put("/", adminHandler.UpdateUser)
			r.Delete("/", adminHandler.DeleteUser)
		})
	})

	r.Route("/problems", func(r chi.Router) {
		r.Post("/", adminHandler.CreateProblem)

		r.Group(func(r chi.Router) {
			r.Use(middleware.OffsetMiddleware)
			r.Use(middleware.LimitMiddleware)
			r.Use(middleware.SearchMiddleware)

			r.Use(middleware.ProblemPremiumFilterMiddleware)
			r.Use(middleware.ProblemPublicFilterMiddleware)
			r.Use(middleware.ProblemSortByFilterMiddleware)
			r.Use(middleware.ProblemSortOrderFilterMiddleware)

			r.Get("/", adminHandler.GetProblems)
		})

		r.Route("/{problemSlug}", func(r chi.Router) {
			r.Use(problemMiddleware)

			r.Get("/", adminHandler.GetProblem)
			r.Put("/", adminHandler.UpdateProblem)
			r.Delete("/", adminHandler.DeleteProblem)
		})
	})

	return r
}
