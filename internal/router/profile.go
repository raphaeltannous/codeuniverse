package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

func profileRouter(
	userHandler *handlers.UserHandler,

	authMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Use(authMiddleware)

	r.Get("/me", handlersutils.Unimplemented)

	return r
}
