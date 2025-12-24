package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `db:"id" json:"-"`

	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	Email        string `db:"email" json:"-"`

	IsVerified bool `db:"is_verified" json:"isVerified"`
	IsActive   bool `db:"is_active" json:"isActive"`

	Role string `db:"role" json:"role"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
