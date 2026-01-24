package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/stripe/stripe-go/v84"
	billingportal "github.com/stripe/stripe-go/v84/billingportal/session"
	"github.com/stripe/stripe-go/v84/checkout/session"
	"github.com/stripe/stripe-go/v84/customer"
)

type StripeService interface {
	GetCustomer(ctx context.Context, user *models.User) (*stripe.Customer, error)
	UpdatePaymentMethod(ctx context.Context, user *models.User) (string, error)
	GetCheckoutSession(ctx context.Context, user *models.User, stripePriceID string) (*stripe.CheckoutSession, error)

	HandleWebhook(ctx context.Context, event stripe.Event)
}

type stripeService struct {
	userRepository repository.UserRepository

	logger *slog.Logger
}

func (s *stripeService) HandleWebhook(ctx context.Context, event stripe.Event) {
	switch event.Type {
	case "invoice.paid":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			s.logger.Error("failed to unmarshal invoice event", "err", err)
			return
		}

		user, err := s.userRepository.GetByStripeCustomerId(
			ctx,
			invoice.Customer.ID,
		)
		if err != nil {
			s.logger.Error("failed to get user via stripe_customer_id", "customerId", invoice.Customer.ID, "err", err)
			return
		}

		err = s.userRepository.UpdatePremiumStatus(
			ctx,
			user.ID,
			"premium",
		)
		if err != nil {
			s.logger.Error("failed to update premium status", "user", user, "newStatus", "premium", "err", err)
			return
		}
	case "invoice.payment_failed":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			s.logger.Error("failed to unmarshal invoice event", "err", err)
			return
		}

		user, err := s.userRepository.GetByStripeCustomerId(
			ctx,
			invoice.Customer.ID,
		)
		if err != nil {
			s.logger.Error("failed to get user via stripe_customer_id", "customerId", invoice.Customer.ID, "err", err)
			return
		}

		err = s.userRepository.UpdatePremiumStatus(
			ctx,
			user.ID,
			"free",
		)
		if err != nil {
			s.logger.Error("failed to update premium status", "user", user, "newStatus", "free", "err", err)
			return
		}
	case "customer.subscription.created", "customer.subscription.deleted", "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			s.logger.Error("failed to unmarshal subscription event", "err", err)
			return
		}

		user, err := s.userRepository.GetByStripeCustomerId(
			ctx,
			sub.Customer.ID,
		)
		if err != nil {
			s.logger.Error("failed to get user via stripe_customer_id", "customerId", sub.Customer.ID, "err", err)
			return
		}

		newStatus := "premium"
		if sub.Status == stripe.SubscriptionStatusCanceled && sub.EndedAt != 0 {
			newStatus = "free"
		}

		err = s.userRepository.UpdatePremiumStatus(
			ctx,
			user.ID,
			newStatus,
		)
		if err != nil {
			s.logger.Error("failed to update premium status for user", "user", user, "newStatus", newStatus, "err", err)
			return
		}
	default:
	}
}

func (s *stripeService) UpdatePaymentMethod(ctx context.Context, user *models.User) (string, error) {
	cust, err := s.GetCustomer(
		ctx,
		user,
	)
	if err != nil {
		return "", err
	}

	portalSession, err := billingportal.New(&stripe.BillingPortalSessionParams{
		Customer:  stripe.String(cust.ID),
		ReturnURL: stripe.String("http://localhost:8080/subscription"),
	})
	if err != nil {
		s.logger.Error("failed to create portal session", "user", user, "err", err)
		return "", err
	}

	return portalSession.URL, nil
}

func (s *stripeService) GetCheckoutSession(
	ctx context.Context,
	user *models.User,
	stripePriceID string,
) (*stripe.CheckoutSession, error) {
	cust, err := s.GetCustomer(
		ctx,
		user,
	)
	if err != nil {
		return nil, err
	}

	params := &stripe.CheckoutSessionParams{
		UIMode:    stripe.String("embedded"),
		ReturnURL: stripe.String("http://localhost:8080/subscription"),
		Customer:  stripe.String(cust.ID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(stripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
	}

	sess, err := session.New(params)
	if err != nil {
		s.logger.Error("failed to create checkout session", "user", user, "err", err)
		return nil, err
	}

	return sess, nil
}

func (s *stripeService) GetCustomer(
	ctx context.Context,
	user *models.User,
) (*stripe.Customer, error) {
	if user.StripeCustomerID != "" {
		cust, err := customer.Get(user.StripeCustomerID, nil)
		if err != nil {
			s.logger.Error("failed to get customer", "user", user, "err", err)
			return nil, err
		}

		return cust, nil
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.Username),
		Metadata: map[string]string{
			"user_id": user.ID.String(),
		},
	}

	cust, err := customer.New(params)
	if err != nil {
		s.logger.Error("failed to create customer", "user", user, "err", err)
		return nil, err
	}

	err = s.userRepository.UpdateStripeCustomerId(ctx, user.ID, cust.ID)
	if err != nil {
		s.logger.Error("failed to update customerId", "user", user, "customerId", cust.ID, "err", err)
		return nil, err
	}

	return cust, nil
}

func NewStripeService(
	userRepository repository.UserRepository,

	key string,
) StripeService {
	stripe.Key = key

	return &stripeService{
		userRepository: userRepository,

		logger: slog.Default().WithGroup("service.StripeService"),
	}
}
