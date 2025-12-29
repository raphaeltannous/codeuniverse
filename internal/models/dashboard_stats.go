package models

import (
	"time"

	"github.com/google/uuid"
)

type DashboardStats struct {
	TotalUsers                 int `json:"totalUsers"`
	TotalUsersRegisteredLast7d int `json:"totalUsersRegisteredLast7d"`

	TotalProblems  int `json:"totalProblems"`
	EasyProblems   int `json:"easyProblems"`
	MediumProblems int `json:"mediumProblems"`
	HardProblems   int `json:"hardProblems"`

	TotalSubmissions   int `json:"totalSubmissions"`
	SubmissionsLast24h int `json:"submissionsLast24h"`
	SubmissionsLast7d  int `json:"submissionsLast7d"`
	SubmissionsLast30d int `json:"submissionsLast30d"`

	TotalAdmins        int `json:"totalAdmins"`
	PendingSubmissions int `json:"pendingSubmissions"`

	AcceptedSubmissions int     `json:"acceptedSubmissions"`
	AcceptanceRate      float64 `json:"acceptanceRate"`
}

type ActivityLog struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Username    string         `json:"username"`
	Description string         `json:"description"`
	Timestamp   time.Time      `json:"timestamp"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type SubmissionActivity struct {
	ID           uuid.UUID
	Username     string
	ProblemTitle string
	ProblemId    uuid.UUID
	Status       string
	CreatedAt    time.Time
}

type DailySubmissions struct {
	Date        string `json:"date"`
	Submissions int    `json:"submissions"`
	Accepted    int    `json:"accepted"`
}

func NewDashboardStats(
	totalUsers int,
	totalUsersRegisteredLast7d int,

	easyProblems int,
	mediumProblems int,
	hardProblems int,

	totalSubmissions int,
	totalSubmissionsLast1d int,
	totalSubmissionsLast7d int,
	totalSubmissionsLast30d int,

	totalAdmins int,

	pendingSubmissions int,
	acceptedSubmissions int,

) *DashboardStats {
	var acceptanceRate float64
	if totalSubmissions > 0 {
		acceptanceRate = (float64(acceptedSubmissions) / float64(totalSubmissions)) * 100
	}

	return &DashboardStats{
		TotalUsers:                 totalUsers,
		TotalUsersRegisteredLast7d: totalUsersRegisteredLast7d,

		TotalProblems:  easyProblems + mediumProblems + hardProblems,
		EasyProblems:   easyProblems,
		MediumProblems: mediumProblems,
		HardProblems:   hardProblems,

		TotalSubmissions:   totalSubmissions,
		SubmissionsLast24h: totalSubmissionsLast1d,
		SubmissionsLast7d:  totalSubmissionsLast7d,
		SubmissionsLast30d: totalSubmissionsLast30d,

		TotalAdmins: totalAdmins,

		PendingSubmissions:  pendingSubmissions,
		AcceptedSubmissions: acceptedSubmissions,

		AcceptanceRate: acceptanceRate,
	}
}
