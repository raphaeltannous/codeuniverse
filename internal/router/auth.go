package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func authRouter(
	userHandler *handlers.UserHandler,

	authMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Route("/signup", func(r chi.Router) {
		r.Post("/", userHandler.Signup)
		r.Post("/verify", userHandler.VerifyEmailByToken)
	})

	r.Route("/login", func(r chi.Router) {
		r.Post("/", userHandler.Login)

		r.Route("/mfa", func(r chi.Router) {
			r.Use(middleware.MfaTokenMiddleware)

			r.Post("/", userHandler.MfaVerification)
			r.Post("/resend", userHandler.ResendMfaVerification)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/logout", userHandler.Logout)
		r.Post("/refresh", userHandler.RefreshJWTToken)
		r.Get("/status", userHandler.JWTTokenStatus)
	})

	r.Route("/password", func(r chi.Router) {
		r.Post("/request", userHandler.PasswordResetRequest)
		r.Post("/reset", userHandler.PasswordResetByToken)
	})

	return r
}
