package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID uuid.UUID `db:"id" json:"-"`

	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`

	Description string `db:"description" json:"description"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}
