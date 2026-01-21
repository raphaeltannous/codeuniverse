package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `db:"id" json:"id"`

	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	Email        string `db:"email" json:"email"`
	AvatarURL    string `db:"avatar_url" json:"avatarUrl"`

	IsVerified bool `db:"is_verified" json:"isVerified"`
	IsActive   bool `db:"is_active" json:"isActive"`

	StripeCustomerID string `db:"stripe_customer_id" json:"stripeCustomerId"`
	PremiumStatus    string `db:"premium_status" json:"premiumStatus"`

	Role string `db:"role" json:"role"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
