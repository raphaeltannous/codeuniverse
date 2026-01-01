package handlers

import (
	"net/http"

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
