package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

const (
	OffsetDefault = 0
	LimitDefault  = 25

	OffsetCtxKey = "offset"
	LimitCtxKey  = "limit"
)

var OffsetMiddleware = makePaginationHandler(
	"offset",
	OffsetCtxKey,
	OffsetDefault,
)

var LimitMiddleware = makePaginationHandler(
	"limit",
	LimitCtxKey,
	LimitDefault,
)

func makePaginationHandler(getParam, ctxKey string, defaultValue int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			number := defaultValue
			if numberString := r.URL.Query().Get(getParam); numberString != "" {
				var err error
				number, err = strconv.Atoi(numberString)

				if err != nil {
					http.Error(
						w,
						fmt.Sprintf(
							"%s should be and integer.",
							getParam,
						),
						http.StatusBadRequest,
					)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ctxKey, number)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
