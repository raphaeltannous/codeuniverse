package handlers

import (
	"net/http"

	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
)

type CourseHandler struct {
	courseService services.CourseService
}

func NewCourseHandler(courseService services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

func (h *CourseHandler) GetPublicCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	courses, err := h.courseService.GetAllCourses(ctx)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := make([]*models.Course, 0, len(courses))
	for _, course := range courses {
		if course.IsPublished {
			response = append(response, course)
		}
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}

func (h *CourseHandler) GetPublicCoursesWithProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	courses, err := h.courseService.GetAllCourses(ctx)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	type courseProgressResponse struct {
		*models.Course
		CompletionPercentage float64 `json:"completionPercentage"`
	}

	response := make([]*courseProgressResponse, 0, len(courses))
	for _, course := range courses {
		if course.IsPublished {
			courseProgress, err := h.courseService.GetCourseProgress(
				ctx,
				course,
				user,
			)
			if err != nil {
				courseProgress = []*models.CourseLessonProgress(nil)
			}

			var completedLessonCount float64
			for _, progress := range courseProgress {
				if progress.IsCompleted {
					completedLessonCount++
				}
			}

			var completionPercentage float64
			if completedLessonCount > 0 {
				completionPercentage = float64(course.TotalLessons) / completionPercentage
			}

			response = append(
				response,
				&courseProgressResponse{
					course,
					completionPercentage,
				},
			)
		}
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)

}

func (h *CourseHandler) GetLessons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	course, ok := ctx.Value(middleware.CourseCtxKey).(*models.Course)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	lessons, err := h.courseService.GetCourseLessons(
		ctx,
		course,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	validLessons := make([]*models.Lesson, 0, len(lessons))
	for _, lesson := range lessons {
		if lesson.DurationSeconds > 0 && lesson.VideoURL != "default.mp4" {
			validLessons = append(validLessons, lesson)
		}
	}

	response := map[string]any{
		"lessons":     validLessons,
		"courseTitle": course.Title,
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}
