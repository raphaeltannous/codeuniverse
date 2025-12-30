package middleware

import (
	"context"
	"errors"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const LessonCtxKey = "lesson"

func LessonMiddleware(next http.Handler, courseService services.CourseService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lessonId := chi.URLParam(r, "lessonId")
		lessonUUID, err := uuid.Parse(lessonId)
		if err != nil {
			handlersutils.WriteResponseJSON(w, handlersutils.NewInternalServerAPIError(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		lesson, err := courseService.GetLesson(
			ctx,
			lessonUUID,
		)
		if err != nil {
			apiError := handlersutils.NewInternalServerAPIError()
			switch {
			case errors.Is(err, repository.ErrLessonNotFound):
				apiError.Code = "LESSON_NOT_FOUND"
				apiError.Message = "Lesson not found."
				handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			default:
				handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
			}
			return
		}

		ctx = context.WithValue(ctx, LessonCtxKey, lesson)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
