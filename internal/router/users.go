package router

import (
	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func usersRouter(
	userHandlers *handlers.UserHandler,
) chi.Router {
	r := chi.NewRouter()

	return r
}
