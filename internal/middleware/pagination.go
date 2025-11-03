package middleware

import (
	"context"
	"net/http"
	"strconv"
)

const (
	OffsetDefault = 0
	LimitDefault  = 25
)

func OffsetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset := OffsetDefault
		if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
			var err error
			offset, err = strconv.Atoi(offsetParam)
			if err != nil {
				http.Error(w, "offset should be an integer", http.StatusBadRequest)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "offset", offset)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := LimitDefault
		if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
			var err error
			limit, err = strconv.Atoi(limitParam)
			if err != nil {
				http.Error(w, "limit should be an integer", http.StatusBadRequest)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "limit", limit)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
