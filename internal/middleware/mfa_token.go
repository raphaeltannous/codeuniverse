package middleware

import (
	"context"
	"net/http"

	"git.riyt.dev/codeuniverse/internal/utils/handlersutils"
)

type MfaRequestBody struct {
	Token string `json:"token"`
	Code  string `json:"code"`
}

func MfaTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody := MfaRequestBody{}

		if !handlersutils.DecodeJSONRequest(w, r, &requestBody) {
			return
		}

		ctx := context.WithValue(r.Context(), "requestBody", requestBody)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
