package middleware

import (
	"context"
	"errors"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/models"
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
				handlersutils.WriteResponseJSON(w, apiError, http.StatusNotFound)
			default:
				handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
			}
			return
		}

		user, ok := ctx.Value(UserAuthCtxKey).(*models.User)
		if !ok {
			user = &models.User{
				Role:          "user",
				PremiumStatus: "free",
			}
		}

		if user.PremiumStatus == "free" && user.Role != "admin" {
			apiError := handlersutils.NewAPIError(
				"COURSE_IS_PREMIUM",
				"Course requires premium plan.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusForbidden)
			return
		}

		if !course.IsPublished && user.Role != "admin" {
			apiError := handlersutils.NewAPIError(
				"COURSE_NOT_FOUND",
				"Course not found.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusNotFound)
			return
		}

		ctx = context.WithValue(ctx, CourseCtxKey, course)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
