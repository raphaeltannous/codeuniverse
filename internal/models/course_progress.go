package models

import (
	"time"

	"github.com/google/uuid"
)

type CourseLessonProgress struct {
	UserId   uuid.UUID `db:"user_id" json:"-"`
	CourseId uuid.UUID `db:"course_id" json:"-"`
	LessonId uuid.UUID `db:"lesson_id" json:"-"`

	IsCompleted bool `db:"is_completed" json:"isCompleted"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
