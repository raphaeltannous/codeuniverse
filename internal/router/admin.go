package router

import (
	"fmt"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func adminRouter(userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.AdminOnly)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	r.Post("/user", userHandler.GetUserInfoById)

	r.Group(func(r chi.Router) {
		r.Use(middleware.OffsetMiddleware)
		r.Use(middleware.LimitMiddleware)

		r.Get("/users", userHandler.GetAllUsers)
	})

	return r
}
