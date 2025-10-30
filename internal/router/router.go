package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Service(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			fmt.Fprintln(w, "not pong")
		} else {
			fmt.Fprintln(w, "pong")
		}
	})

	return r
}
