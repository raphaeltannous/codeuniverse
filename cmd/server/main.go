package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"git.riyt.dev/codeuniverse/internal/database"
	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/judger"
	"git.riyt.dev/codeuniverse/internal/logger"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/repository/postgres"
	"git.riyt.dev/codeuniverse/internal/router"
	"git.riyt.dev/codeuniverse/internal/services"
)

func main() {
	// TODO: command line option for logging level
	lg, err := logger.New(slog.LevelDebug)
	if err != nil {
		log.Fatal(err)
	}
	slog.SetDefault(lg)

	judge, err := judger.NewJudge()
	if err != nil {
		log.Fatal(err)
	}
	defer judge.Cli.Close()

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

	if err := judge.InitializeContainers(ctx); err != nil {
		log.Fatal(err)
	}

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
	//
	// middlewares
	authMiddleware := func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, userService)
	}

	return router.Service(userHandler, authMiddleware)
}
