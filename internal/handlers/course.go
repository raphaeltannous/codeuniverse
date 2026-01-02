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

	courses, err := h.courseService.GetAllPublicCourses(ctx)
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

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	courses, err := h.courseService.GetAllPublicCourses(ctx)
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
				completionPercentage = (completedLessonCount / float64(course.TotalLessons)) * 100
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

func (h *CourseHandler) UpdateIsCompleted(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	course, ok := ctx.Value(middleware.CourseCtxKey).(*models.Course)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	lesson, ok := ctx.Value(middleware.LessonCtxKey).(*models.Lesson)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	lessonProgress := &models.CourseLessonProgress{
		UserId:      user.ID,
		CourseId:    course.ID,
		LessonId:    lesson.ID,
		IsCompleted: true,
	}
	err := h.courseService.UpdateCourseProgress(
		ctx,
		lessonProgress,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CourseHandler) GetCourseProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserAuthCtxKey).(*models.User)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	course, ok := ctx.Value(middleware.CourseCtxKey).(*models.Course)
	if !ok {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	courseProgress, err := h.courseService.GetCourseProgress(
		ctx,
		course,
		user,
	)
	if err != nil {
		handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
		return
	}

	response := make(map[string]bool)
	for _, progress := range courseProgress {
		response[progress.LessonId.String()] = progress.IsCompleted
	}

	handlersutils.WriteResponseJSON(w, response, http.StatusOK)
}
