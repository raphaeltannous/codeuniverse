package router

import (
	"net/http"
	"os"
	"path/filepath"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Service(
	userHandler *handlers.UserHandler,
	problemsHandler *handlers.ProblemHandler,
	statsHandler *handlers.StatsHandler,
	staticHandler *handlers.StaticHandler,
	adminHandler *handlers.AdminHandler,
	courseHandler *handlers.CourseHandler,

	authMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
	courseMiddleware func(next http.Handler) http.Handler,
	lessonMiddleware func(next http.Handler) http.Handler,
	userMiddleware func(next http.Handler) http.Handler,
	hintMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	// r.Use(func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		time.Sleep(1 * time.Second)
	// 		next.ServeHTTP(w, r)
	// 	})
	// })

	r.Mount("/api", apiRouter(
		userHandler,
		problemsHandler,
		statsHandler,
		staticHandler,
		adminHandler,
		courseHandler,

		authMiddleware,
		problemMiddleware,
		courseMiddleware,
		lessonMiddleware,
		userMiddleware,
		hintMiddleware,
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
	courseHandler *handlers.CourseHandler,

	authMiddleware func(next http.Handler) http.Handler,
	problemMiddleware func(next http.Handler) http.Handler,
	courseMiddleware func(next http.Handler) http.Handler,
	lessonMiddleware func(next http.Handler) http.Handler,
	userMiddleware func(next http.Handler) http.Handler,
	hintMiddleware func(next http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Mount("/health", heathRouter())
	r.Mount("/auth", authRouter(userHandler, authMiddleware))

	r.Mount("/problems", problemsRouter(problemsHandler, authMiddleware, problemMiddleware))
	r.Mount("/courses", courseRouter(courseHandler, authMiddleware, courseMiddleware, lessonMiddleware))

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
		userMiddleware,
		problemMiddleware,
		hintMiddleware,
	))

	r.Mount("/static", staticRouter(
		staticHandler,
		authMiddleware,
		lessonMiddleware,
	))

	return r
}
