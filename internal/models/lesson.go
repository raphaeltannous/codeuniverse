package models

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	ID uuid.UUID `db:"id" json:"id"`

	CourseId uuid.UUID `db:"course_id" json:"courseId"`

	Title           string `db:"title" json:"title"`
	Description     string `db:"description" json:"description"`
	VideoURL        string `db:"video_url" json:"videoURL"`
	DurationSeconds int    `db:"duration_seconds" json:"durationSeconds"`
	LessonNumber    int    `db:"lesson_number" json:"lessonNumber"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
