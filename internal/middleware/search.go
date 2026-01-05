package middleware

import (
	"context"
	"net/http"
)

const (
	SearchCtxKey = "search"
)

func SearchMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if searchParam := r.URL.Query().Get(SearchCtxKey); searchParam != "" {
			ctx := context.WithValue(r.Context(), SearchCtxKey, searchParam)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
