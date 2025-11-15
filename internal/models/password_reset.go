package models

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	ID uuid.UUID `db:"id" json:"-"`

	UserId uuid.UUID `db:"user_id" json:"-"`
	Hash   string    `db:"token_hash" json:"-"`

	ExpiresAt time.Time `db:"expires_at" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

func (pr PasswordReset) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ID", pr.ID.String()),
		slog.String("UserId", pr.UserId.String()),
		slog.String("Hash", pr.Hash),
		slog.Time("ExpiresAt", pr.ExpiresAt),
		slog.Time("CreatedAt", pr.CreatedAt),
	)
}
