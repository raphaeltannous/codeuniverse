package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func heathRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	return r
}
