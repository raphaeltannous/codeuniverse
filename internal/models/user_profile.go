package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID uuid.UUID `db:"user_id" json:"-"`

	Name              *string `db:"name" json:"name"`
	Bio               *string `db:"bio" json:"bio"`
	AvatarURL         *string `db:"avatar_url" json:"avatarUrl"`
	Country           *string `db:"country" json:"country"`
	PreferredLanguage *string `db:"preferred_language" json:"preferredLanguage"`

	TotalSubmissions    *int `db:"total_submissions" json:"totalSubmissions"`
	AcceptedSubmissions *int `db:"accepted_submissions" json:"acceptedSubmissions"`
	ProblemsSolved      *int `db:"problems_solved" json:"problemsSolved"`
	EasySolved          *int `db:"easy_solved" json:"easySolved"`
	MediumSolved        *int `db:"medium_solved" json:"mediumSolved"`
	HardSolved          *int `db:"hard_solved" json:"hardSolved"`

	WebsiteURL  *string `db:"website_url" json:"websiteUrl"`
	GithubURL   *string `db:"github_url" json:"githubUrl"`
	LinkedinURL *string `db:"linkedin_url" json:"linkedinUrl"`
	XURL        *string `db:"x_url" json:"xUrl"`

	LastActive *time.Time `db:"last_active" json:"lastActive"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
