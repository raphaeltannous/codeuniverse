package middleware

import (
	"context"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/services"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const ProblemHintCtxKey = "problemHint"

func ProblemHintMiddleware(next http.Handler, problemService services.ProblemService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hintId := chi.URLParam(r, "hintId")
		hintUUId, err := uuid.Parse(hintId)
		if err != nil {
			apiError := handlersutils.NewAPIError(
				"INVALID_HINT_ID",
				"Invalid hint id.",
			)
			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		hint, err := problemService.GetHint(
			ctx,
			hintUUId,
		)
		if err != nil {
			apiError := handlersutils.NewInternalServerAPIError()
			switch err {
			case repository.ErrProblemHintNotFound:
				apiError.Code = "HINT_NOT_FOUND"
				apiError.Message = "Hint not found."

				handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			default:
				handlersutils.WriteResponseJSON(w, apiError, http.StatusInternalServerError)
			}
			return
		}

		ctx = context.WithValue(ctx, ProblemHintCtxKey, hint)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
