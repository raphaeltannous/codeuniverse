package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"git.riyt.dev/codeuniverse/internal/database"
	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/repository/postgres"
	"git.riyt.dev/codeuniverse/internal/router"
	"git.riyt.dev/codeuniverse/internal/services"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := &http.Server{
		Addr:    ":3333",
		Handler: service(db),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}

func service(db *sql.DB) http.Handler {
	// repos
	userRepo := postgres.NewUserRepository(db)
	//
	// services
	userService := services.NewUserService(userRepo)
	//
	// handlers
	userHandler := handlers.NewUserHandler(userService)

	return router.Service(userHandler)
}
