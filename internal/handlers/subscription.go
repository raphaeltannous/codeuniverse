package handlers

import (
	"io"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/stripe/stripe-go/v84/webhook"
)

type SubscriptionHandler struct {
	stripeService services.StripeService
}

func NewSubscriptionHandler(
	stripeService services.StripeService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		stripeService: stripeService,
	}
}

func (h *SubscriptionHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	if user.StripeCustomerID == "" {
		cust, err := h.stripeService.GetCustomer(ctx, user)
		if err != nil {
			handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
			return
		}
		user.StripeCustomerID = cust.ID
	}

	response := map[string]string{
		"customerId": user.StripeCustomerID,
		"email":      user.ID.String(),
		"status":     user.PremiumStatus,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *SubscriptionHandler) GetCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		PriceId string `json:"priceId"`
	}

	if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
		return
	}

	ctx := r.Context()
	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	sess, err := h.stripeService.GetCheckoutSession(
		ctx,
		user,
		requestBody.PriceId,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":      "Checkout session created.",
		"clientSecret": sess.ClientSecret,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *SubscriptionHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	portalUrl, err := h.stripeService.UpdatePaymentMethod(
		ctx,
		user,
	)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"FAILED_TO_GET_PORTAL",
			"Failed to get portal. Please contact support.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":   "Redirecting to payment portal...",
		"portalUrl": portalUrl,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *SubscriptionHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), "whsec_2036c4e5a8c105044b7d676b5982cca2f2c79d5393550f62912669fe14f1d95e")
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	h.stripeService.HandleWebhook(ctx, event)
}
