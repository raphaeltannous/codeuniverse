package middleware

import (
	"fmt"
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("admin middleware")
		next.ServeHTTP(w, r)
	})
}
