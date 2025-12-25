package handlers

import "git.riyt.dev/codeuniverse/internal/services"

type CourseHandler struct {
	courseService services.CourseService
}

func NewCourseHandler(courseService services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}
