package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Service(
	userHandler *handlers.UserHandler,
	problemsHandler *handlers.ProblemHanlder,

	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)

	r.Mount("/api", apiRouter(
		userHandler,
		problemsHandler,

		authMiddleware,
	))

	return r
}

func apiRouter(
	userHandler *handlers.UserHandler,
	problemsHandler *handlers.ProblemHanlder,

	authMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Mount("/health", heathRouter())
	r.Mount("/auth", authRouter(userHandler, authMiddleware))

	r.Mount("/problems", problemsRouter(problemsHandler, authMiddleware))
	r.Mount("/submissions", submissionsRouter(authMiddleware))

	r.Mount("/admin", adminRouter(userHandler, problemsHandler, authMiddleware))

	return r
}
