package middleware

import (
	"context"
	"errors"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

const CourseCtxKey = "course"

func CourseMiddleware(next http.Handler, courseService services.CourseService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		courseSlug := chi.URLParam(r, "courseSlug")
		ctx := r.Context()

		course, err := courseService.GetCourseBySlug(
			ctx,
			courseSlug,
		)
		if err != nil {
			apiError := handlersutils.NewInternalServerAPIError()
			switch {
			case errors.Is(err, repository.ErrCourseNotFound):
				apiError.Code = "COURSE_NOT_FOUND"
				apiError.Message = "Course not found."
				handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			default:
				handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
			}
			return
		}

		ctx = context.WithValue(ctx, CourseCtxKey, course)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
