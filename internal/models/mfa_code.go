package models

import (
	"time"

	"github.com/google/uuid"
)

type MfaCode struct {
	ID uuid.UUID `db:"id" json:"-"`

	UserId    uuid.UUID `db:"user_id" json:"-"`
	TokenHash string    `db:"token_hash" json:"-"`
	CodeHash  string    `db:"code_hash" json:"-"`

	ExpiresAt time.Time `db:"expires_at" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}
