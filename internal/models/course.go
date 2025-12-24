package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID uuid.UUID `db:"id" json:"id"`

	Title        string `db:"title" json:"title"`
	Description  string `db:"title" json:"description"`
	ThumbnailURL string `db:"thumbnail_url" json:"thumbnailURL"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
