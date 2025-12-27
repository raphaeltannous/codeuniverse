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
	staticHandler *handlers.StaticHandler,

	authMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)

	r.Mount("/api", apiRouter(
		userHandler,
		problemsHandler,
		staticHandler,

		authMiddleware,
		problemMiddleware,
	))

	return r
}

func apiRouter(
	userHandler *handlers.UserHandler,
	problemsHandler *handlers.ProblemHandler,
	staticHandler *handlers.StaticHandler,

	authMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Mount("/health", heathRouter())
	r.Mount("/auth", authRouter(userHandler, authMiddleware))

	r.Mount("/problems", problemsRouter(problemsHandler, authMiddleware, problemMiddleware))
	r.Mount("/submissions", submissionsRouter(authMiddleware))

	r.Mount("/users", usersRouter(userHandler))
	r.Mount("/profile", profileRouter(userHandler, authMiddleware))
	r.Mount("/admin", adminRouter(userHandler, problemsHandler, authMiddleware))

	r.Mount("/static", staticRouter(staticHandler, authMiddleware))

	return r
}
