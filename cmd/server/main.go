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
	defer judge.Close()

	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := &http.Server{
		Addr:    ":3333",
		Handler: service(db, mailMan, *judge),
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

func service(
	db *sql.DB,
	mailMan mailer.Mailer,
	judge judger.Judge,
) http.Handler {
	// repos
	userRepo := postgres.NewUserRepository(db)
	userProfileRepo := postgres.NewUserProfileRepository(db)
	problemRepository := postgres.NewProblemRepository(db)
	problemNoteRepository := postgres.NewProblemNoteRepository(db)
	runRepository := postgres.NewRunRepository(db)
	submissionRepository := postgres.NewSubmissionRepository(db)
	mfaRepo := postgres.NewMfaCodeRepository(db)
	passwordResetRepo := postgres.NewPasswordResetRepository(db)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(db)

	dbTransactor := postgres.NewPostgreSQLTransactor(db)

	// services
	userService := services.NewUserService(
		userRepo,
		userProfileRepo,
		submissionRepository,
		problemRepository,
		mfaRepo,
		passwordResetRepo,
		emailVerificationRepo,
		dbTransactor,

		mailMan,
	)

	problemService := services.NewProblemService(
		problemRepository,
		problemNoteRepository,
		runRepository,
		submissionRepository,

		judge,
	)

	staticService := services.NewStaticService()

	// handlers
	userHandler := handlers.NewUserHandler(userService)
	problemHandler := handlers.NewProblemsHandlers(problemService)
	staticHandler := handlers.NewStaticHandler(staticService)

	// middlewares
	authMiddleware := func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, userService)
	}
	problemMiddleware := func(next http.Handler) http.Handler {
		return middleware.ProblemMiddleware(next, problemService)
	}

	return router.Service(
		userHandler,
		problemHandler,
		staticHandler,

		authMiddleware,
		problemMiddleware,
	)
}
