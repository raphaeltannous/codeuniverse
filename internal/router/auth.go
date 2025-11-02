package router

import (
	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func authRouter(userHandler *handlers.UserHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", userHandler.CreateUser)

	return r
}
