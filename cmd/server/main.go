package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.riyt.dev/codeuniverse/internal/database"
	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/judger"
	"git.riyt.dev/codeuniverse/internal/logger"
	"git.riyt.dev/codeuniverse/internal/mailer"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/repository/postgres"
	"git.riyt.dev/codeuniverse/internal/router"
	"git.riyt.dev/codeuniverse/internal/services"
)

var codeuniverseEnv string

func init() {
	allowedEnv := map[string]bool{
		"production":  true,
		"development": true,
	}

	codeuniverseEnv = os.Getenv("CODEUNIVERSE_ENV")
	if codeuniverseEnv == "" {
		log.Fatal("CODEUNIVERSE_ENV is not set")
	} else if !allowedEnv[codeuniverseEnv] {
		log.Fatal("CODEUNIVERSE_ENV should either be production, or development")
	}
}

func main() {
	var mailMan mailer.Mailer

	switch codeuniverseEnv {
	case "development":
		mailMan = mailer.NewSMTPMailer(
			"localhost",
			1025,
			"codeuniverse.lb@gmail.com",
			"",
			"codeuniverse.lb@gmail.com",
		)
	case "production":
		gmailSMTPPassword := os.Getenv("CODEUNIVERSE_SMTP_GMAIL_PASSWORD")
		if gmailSMTPPassword == "" {
			log.Fatal("CODEUNIVERSE_SMTP_GMAIL_PASSWORD is not set")
		}
		mailMan = mailer.NewSMTPMailer(
			"smtp.gmail.com",
			587,
			"codeuniverse.lb@gmail.com",
			gmailSMTPPassword,
			"codeuniverse.lb@gmail.com",
		)
	}
	mailMan.Send(context.Background(), "log@riyt.dev", "Testing Mailpit", "Hello mailpit")

	// TODO: command line option for logging level or using environment variables
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
		Handler: service(db, mailMan),
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

func service(db *sql.DB, mailMan mailer.Mailer) http.Handler {
	// repos
	userRepo := postgres.NewUserRepository(db)
	//
	// services
	userService := services.NewUserService(userRepo, mailMan)
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
