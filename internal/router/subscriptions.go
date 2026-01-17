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

	r.Use(authMiddleware)

	r.Get("/status", subscriptionHandler.GetStatus)
	r.Post("/checkout-session", subscriptionHandler.GetCheckoutSession)
	r.Post("/cancel", subscriptionHandler.CancelSubscription)
	r.Post("/update-payment-method", subscriptionHandler.UpdatePaymentMethod)

	return r
}
