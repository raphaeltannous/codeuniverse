package services

import (
	"context"
	"log/slog"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/stripe/stripe-go/v84"
	billingportal "github.com/stripe/stripe-go/v84/billingportal/session"
	"github.com/stripe/stripe-go/v84/checkout/session"
	"github.com/stripe/stripe-go/v84/customer"
	"github.com/stripe/stripe-go/v84/subscription"
)

type StripeService interface {
	GetCustomer(ctx context.Context, user *models.User) (*stripe.Customer, error)
	UpdatePaymentMethod(ctx context.Context, user *models.User) (string, error)
	CancelSubscription(ctx context.Context, user *models.User) error
	GetSubscriptionStatus(ctx context.Context, user *models.User) (string, error)
	GetCheckoutSession(ctx context.Context, user *models.User, stripePriceID string) (*stripe.CheckoutSession, error)
}

type stripeService struct {
	userRepository repository.UserRepository

	logger *slog.Logger
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

func (s *stripeService) CancelSubscription(ctx context.Context, user *models.User) error {
	cust, err := s.GetCustomer(
		ctx,
		user,
	)
	if err != nil {
		return err
	}

	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(cust.ID),
	}

	subList := subscription.List(params)
	if subList.Err() != nil {
		return subList.Err()
	}

	for subList.Next() {
		sub := subList.Subscription()

		if sub.Status == stripe.SubscriptionStatusActive {
			_, err := subscription.Cancel(sub.ID, nil)
			if err != nil {
				s.logger.Error("failed to cancel subscription")
				return err
			}
		}
	}

	return nil
}

func (s *stripeService) GetSubscriptionStatus(ctx context.Context, user *models.User) (string, error) {
	cust, err := s.GetCustomer(
		ctx,
		user,
	)
	if err != nil {
		return "", err
	}

	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(cust.ID),
	}

	subList := subscription.List(params)
	if subList.Err() != nil {
		return "free", subList.Err()
	}

	for subList.Next() {
		sub := subList.Subscription()
		s.logger.Debug("status", "sub", sub)
		if sub.Status == stripe.SubscriptionStatusActive {
			return "premium", nil
		}
	}

	return "free", nil
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
		ReturnURL: stripe.String("http://localhost:8080/subscriptions"),
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
		logger:         slog.Default().WithGroup("service.StripeService"),
	}
}
