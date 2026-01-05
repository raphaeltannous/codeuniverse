package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
)

func makeRepositoryParamMiddleware[V repository.UserParam | repository.ProblemParam | models.ProblemDifficulty](
	getParam, ctxKey string,
	allowedFilter map[string]V,
) func(next http.Handler) http.Handler {
	allowedFilter[""] = 0

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			filter := r.URL.Query().Get(getParam)
			if _, ok := allowedFilter[filter]; !ok {
				apiError := handlersutils.NewAPIError(
					fmt.Sprintf("UNALLOWED_%s_FILTER", strings.ToUpper(getParam)),
					fmt.Sprintf("The requested %s filter is not allowed.", getParam),
				)

				handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				ctxKey,
				allowedFilter[filter],
			)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
