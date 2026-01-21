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

const (
	ProblemCtxKey = "problem"

	ProblemPublicFilterCtxKey     = "problemPublicFilter"
	ProblemPremiumFilterCtxKey    = "problemPremiumFilter"
	ProblemDifficultyFilterCtxKey = "problemDifficultyFilter"
	ProblemSortByFilterCtxKey     = "problemSortByFilter"
	ProblemSortOrderFilterCtxKey  = "problemSortOrderFilter"
)

var ProblemPremiumFilterMiddleware = makeRepositoryParamMiddleware(
	"premium",
	ProblemPremiumFilterCtxKey,
	map[string]repository.ProblemParam{
		"premium": repository.ProblemPremium,
		"free":    repository.ProblemFree,
	},
)

var ProblemPublicFilterMiddleware = makeRepositoryParamMiddleware(
	"public",
	ProblemPublicFilterCtxKey,
	map[string]repository.ProblemParam{
		"public":  repository.ProblemPublic,
		"private": repository.ProblemPrivate,
	},
)

var ProblemDifficultyFilterMiddleware = makeRepositoryParamMiddleware(
	"difficulty",
	ProblemDifficultyFilterCtxKey,
	map[string]models.ProblemDifficulty{
		"easy":   models.ProblemEasy,
		"medium": models.ProblemMedium,
		"hard":   models.ProblemHard,
	},
)

var ProblemSortByFilterMiddleware = makeRepositoryParamMiddleware(
	"sortBy",
	ProblemSortByFilterCtxKey,
	map[string]repository.ProblemParam{
		"title":     repository.ProblemSortByTitle,
		"createdAt": repository.ProblemSortByCreatedAt,
	},
)

var ProblemSortOrderFilterMiddleware = makeRepositoryParamMiddleware(
	"sortOrder",
	ProblemSortOrderFilterCtxKey,
	map[string]repository.ProblemParam{
		"asc":  repository.ProblemSortOrderAsc,
		"desc": repository.ProblemSortOrderDesc,
	},
)

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

		user, ok := ctx.Value(UserAuthCtxKey).(*models.User)
		if !ok {
			user = &models.User{
				Role:          "user",
				PremiumStatus: "free",
			}
		}

		if problem.IsPremium && user.PremiumStatus == "free" && user.Role != "admin" {
			apiError := handlersutils.NewAPIError(
				"PROBLEM_IS_PREMIUM",
				"Problem requires premium plan.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, ProblemCtxKey, problem)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
