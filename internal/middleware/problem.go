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

const ProblemCtxKey = "problem"

func ProblemMiddleware(next http.Handler, problemService services.ProblemService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		problemSlug := chi.URLParam(r, "problemSlug")
		ctx := r.Context()

		problem, err := problemService.GetBySlug(
			ctx,
			problemSlug,
		)
		if err != nil {
			apiError := handlersutils.NewInternalServerAPIError()
			switch {
			case errors.Is(err, repository.ErrProblemNotFound):
				apiError.Code = "PROBLEM_NOT_FOUND"
				apiError.Message = "Problem not found."
				handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			default:
				handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
			}
			return
		}

		ctx = context.WithValue(ctx, ProblemCtxKey, problem)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
