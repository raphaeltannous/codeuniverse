package models

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID uuid.UUID `db:"id" json:"-"`

	UserId uuid.UUID `db:"user_id" json:"-"`
	Hash   string    `db:"token_hash" json:"-"`

	ExpiresAt time.Time `db:"expires_at" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

func (ev EmailVerification) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("ID", ev.ID.String()),

		slog.String("UserId", ev.UserId.String()),
		slog.String("Hash", ev.Hash),

		slog.Time("ExpiresAt", ev.ExpiresAt),
		slog.Time("CreatedAt", ev.CreatedAt),
	)
}
