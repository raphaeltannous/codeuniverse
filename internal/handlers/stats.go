package handlers

import (
	"fmt"
	"net/http"
	"slices"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
)

type StatsHandler struct {
	userService    services.UserService
	problemService services.ProblemService
}

func NewStatsHandler(
	userService services.UserService,
	problemService services.ProblemService,
) *StatsHandler {
	return &StatsHandler{
		userService:    userService,
		problemService: problemService,
	}
}

func (h *StatsHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalUsers, err := h.userService.GetUsersCount(ctx)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	totalUsersRegisteredLast7d, err := h.userService.GetUsersRegisteredLastNDaysCount(ctx, 7)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	easyCount, err := h.problemService.GetCount(ctx, models.ProblemEasy)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	mediumCount, err := h.problemService.GetCount(ctx, models.ProblemMedium)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	hardCount, err := h.problemService.GetCount(ctx, models.ProblemHard)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	totalSubmissions, err := h.problemService.GetSubmissionsCount(ctx)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	totalSubmissionsLast1d, err := h.problemService.GetSubmissionsLastNDaysCount(ctx, 1)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	totalSubmissionsLast7d, err := h.problemService.GetSubmissionsLastNDaysCount(ctx, 7)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	totalSubmissionsLast30d, err := h.problemService.GetSubmissionsLastNDaysCount(ctx, 30)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	totalAdmins, err := h.userService.GetAdminCount(ctx)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	pendingSubmissions, err := h.problemService.GetPendingSubmissionsCount(ctx)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	acceptedSubmissions, err := h.problemService.GetAcceptedSubmissionsCount(ctx)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	stats := models.NewDashboardStats(
		totalUsers,
		totalUsersRegisteredLast7d,

		easyCount,
		mediumCount,
		hardCount,

		totalSubmissions,
		totalSubmissionsLast1d,
		totalSubmissionsLast7d,
		totalSubmissionsLast30d,

		totalAdmins,

		pendingSubmissions,
		acceptedSubmissions,
	)

	handlersutils.WriteResponseJSON(w, stats, http.StatusOK)
}

func (h *StatsHandler) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	submissionsActivities, err := h.problemService.GetRecentSubmissions(ctx, 10)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	newUsers, err := h.userService.GetRecentRegisteredUsers(ctx, 10)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	activities := make([]*models.ActivityLog, 0, 20)

	for _, sub := range submissionsActivities {
		activities = append(activities, &models.ActivityLog{
			ID:          sub.ID.String(),
			Type:        "submission",
			Username:    sub.Username,
			Description: fmt.Sprintf("Submitted solution for %s", sub.ProblemTitle),
			Timestamp:   sub.CreatedAt,
			Metadata: map[string]any{
				"problem_id":    sub.ProblemId,
				"problem_title": sub.ProblemTitle,
				"status":        sub.Status,
			},
		})
	}

	for _, user := range newUsers {
		activities = append(activities, &models.ActivityLog{
			ID:          user.ID.String(),
			Type:        "user_registration",
			Username:    user.Username,
			Description: "Registered new account",
			Timestamp:   user.CreatedAt,
			Metadata: map[string]any{
				"role": user.Role,
			},
		})
	}

	slices.SortFunc(activities, func(a, b *models.ActivityLog) int {
		return a.Timestamp.Compare(b.Timestamp)
	})

	handlersutils.WriteResponseJSON(w, activities, http.StatusOK)
}

func (h *StatsHandler) GetSubmissionTrendsSample(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rangeParam := r.URL.Query().Get("range")
	var since int
	switch rangeParam {
	case "24h":
		since = 1
	case "7d":
		since = 7
	case "30d":
		since = 30
	default:
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	dailySubmissions, err := h.problemService.GetDailySubmissions(ctx, since)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	handlersutils.WriteResponseJSON(w, dailySubmissions, http.StatusOK)
}
