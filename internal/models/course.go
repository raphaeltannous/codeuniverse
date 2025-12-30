package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID uuid.UUID `db:"id" json:"id"`

	Title        string `db:"title" json:"title"`
	Slug         string `db:"slug" json:"slug"`
	Description  string `db:"title" json:"description"`
	ThumbnailURL string `db:"thumbnail_url" json:"thumbnailURL"`
	Difficulty   string `db:"difficulty" json:"difficulty"`
	IsPublished  bool   `db:"is_published" json:"isPublished"`

	TotalLessons int `json:"totalLessons"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
