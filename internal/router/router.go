package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Service(
	userHandler *handlers.UserHandler,
	problemsHandler *handlers.ProblemHandler,
	statsHandler *handlers.StatsHandler,
	staticHandler *handlers.StaticHandler,
	adminHandler *handlers.AdminHandler,

	authMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
	courseMiddleware func(next http.Handler) http.Handler,
	lessonMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)

	r.Mount("/api", apiRouter(
		userHandler,
		problemsHandler,
		statsHandler,
		staticHandler,
		adminHandler,

		authMiddleware,
		problemMiddleware,
		courseMiddleware,
		lessonMiddleware,
	))

	frontendDir := "./dist"
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path

		wanted := filepath.Join(frontendDir, urlPath)
		if _, err := os.Stat(wanted); err == nil {
			http.ServeFile(w, r, wanted)
			return
		}

		http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
	})

	return r
}

func apiRouter(
	userHandler *handlers.UserHandler,
	problemsHandler *handlers.ProblemHandler,
	statsHandler *handlers.StatsHandler,
	staticHandler *handlers.StaticHandler,
	adminHandler *handlers.AdminHandler,

	authMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
	courseMiddleware func(next http.Handler) http.Handler,
	lessonMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Mount("/health", heathRouter())
	r.Mount("/auth", authRouter(userHandler, authMiddleware))

	r.Mount("/problems", problemsRouter(problemsHandler, authMiddleware, problemMiddleware))
	r.Mount("/submissions", submissionsRouter(authMiddleware))

	r.Mount("/users", usersRouter(userHandler))
	r.Mount("/profile", profileRouter(userHandler, authMiddleware))
	r.Mount("/admin", adminRouter(
		userHandler,
		problemsHandler,
		statsHandler,
		adminHandler,

		authMiddleware,
		courseMiddleware,
		lessonMiddleware,
	))

	r.Mount("/static", staticRouter(
		staticHandler,
		authMiddleware,
		lessonMiddleware,
	))

	return r
}
