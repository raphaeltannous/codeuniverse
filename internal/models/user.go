package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Username      string    `db:"username" json:"username"`
	PasswordHash  string    `db:"password_hash" json:"-"`
	Email         string    `db:"email" json:"email"`
	EmailVerified bool      `db:"email_verified" json:"emailVerified"`
	IsActive      bool      `db:"is_active" json:"isActive"`
	Role          string    `db:"role" json:"role"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}
