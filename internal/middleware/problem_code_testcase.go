package middleware

import (
	"context"
	"net/http"
	"strconv"

	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
	"github.com/go-chi/chi/v5"
)

const ProblemCodeTestcaseIdCtxKey = "problemCodeTestcaseId"

func ProblemCodeTestcaseIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testcaseIdParam := chi.URLParam(r, "testcaseId")

		testcaseId, err := strconv.Atoi(testcaseIdParam)
		if err != nil {
			apiError := handlersutils.NewAPIError(
				"INVALID_TESTCASE_ID",
				"Testcase id should be integer.",
			)

			handlersutils.WriteResponseJSON(w, apiError, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ProblemCodeTestcaseIdCtxKey, testcaseId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
