package models

import (
	"log/slog"
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

func (mc MfaCode) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ID", mc.ID.String()),

		slog.String("UserId", mc.UserId.String()),
		slog.String("TokenHash", mc.TokenHash),
		slog.String("CodeHash", mc.CodeHash),

		slog.Time("ExpiresAt", mc.ExpiresAt),
		slog.Time("CreatedAt", mc.CreatedAt),
	)
}
