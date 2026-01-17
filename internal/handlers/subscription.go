package handlers

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
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

	cust, err := h.stripeService.GetCustomer(ctx, user)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}
	user.StripeCustomerID = cust.ID

	status, err := h.stripeService.GetSubscriptionStatus(ctx, user)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"customerId": cust.ID,
		"email":      user.ID.String(),
		"status":     status,
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

func (h *SubscriptionHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	err := h.stripeService.CancelSubscription(
		ctx,
		user,
	)
	if err != nil {
		apiError := handlersutils.NewAPIError(
			"FAILED_TO_CANCEL_SUBSCRIPTION",
			"Failed to cancel subscription. Please contact support.",
		)

		handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
		return
	}

	handlersutils.WriteSuccessMessage(
		w,
		"Subscription cancelled.",
		http.StatusOK,
	)
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
