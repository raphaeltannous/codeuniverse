package middleware

import (
	"context"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

const ProblemLanguageCtxKey = "problemLanguage"

func ProblemLanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		languageSlugParam := chi.URLParam(r, "languageSlug")

		language, err := models.NewProblemLanguage(languageSlugParam)
		if err != nil {
			apiError := handlersutils.NewAPIError(
				"INVALID_LANGUAGE",
				"Invalid language",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ProblemLanguageCtxKey, language)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
