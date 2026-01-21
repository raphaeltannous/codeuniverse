package router

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func subscriptionRouter(
	subscriptionHandler *handlers.SubscriptionHandler,

	authMiddleware func(next http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	r.Post("/webhook", subscriptionHandler.StripeWebhook)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/status", subscriptionHandler.GetStatus)
		r.Post("/checkout-session", subscriptionHandler.GetCheckoutSession)
		r.Post("/update-payment-method", subscriptionHandler.UpdatePaymentMethod)
	})

	return r
}
