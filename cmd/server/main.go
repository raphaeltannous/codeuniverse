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
	"path/filepath"
	"syscall"
	"time"

	"git.riyt.dev/codeuniverse/internal/database"
	"git.riyt.dev/codeuniverse/internal/handlers"
	"git.riyt.dev/codeuniverse/internal/judger"
	"git.riyt.dev/codeuniverse/internal/logger"
	"git.riyt.dev/codeuniverse/internal/mailer"
	"git.riyt.dev/codeuniverse/internal/middleware"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository/filesystem"
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

	if codeuniverseEnv == "production" {
		models.Domain = "https://codeuniverse.riyt.dev"
	} else {
		models.Domain = "http://localhost:8080"
	}
}

var (
	problemsDataDir string
	stripeSecret    string
)

func init() {
	problemsDataDir = os.Getenv("CODEUNIVERSE_PROBLEMS_DATA_DIR")

	if problemsDataDir == "" {
		log.Fatal("CODEUNIVERSE_PROBLEMS_DATA_DIR is not set.")
	}

	absPath, err := filepath.Abs(problemsDataDir)
	if err != nil {
		log.Fatal("failed to convert CODEUNIVERSE_PROBLEMS_DATA_DIR to absolute path.")
	}

	problemsDataDir = absPath
	slog.Info("problemsDataDir is updated.", "problemsDataDir", problemsDataDir)

	stripeSecret = os.Getenv("CODEUNIVERSE_STRIPE_DEV_KEY")

	if stripeSecret == "" {
		log.Fatal("CODEUNIVERSE_STRIPE_DEV_KEY is not set.")
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
	// repositories
	userRepository := postgres.NewUserRepository(db)
	userProfileRepository := postgres.NewUserProfileRepository(db)
	problemRepository := postgres.NewProblemRepository(db)
	problemNoteRepository := postgres.NewProblemNoteRepository(db)
	problemHintRepository := postgres.NewProblemHintRepository(db)
	courseRepository := postgres.NewPostgresCourseRepository(db)
	lessonRepository := postgres.NewPostgresLessonRepository(db)
	courseProgressRepository := postgres.NewCourseProgressRepository(db)
	runRepository := postgres.NewRunRepository(db)
	submissionRepository := postgres.NewSubmissionRepository(db)
	mfaRepository := postgres.NewMfaCodeRepository(db)
	passwordResetRepository := postgres.NewPasswordResetRepository(db)
	emailVerificationRepository := postgres.NewEmailVerificationRepository(db)

	problemCodeRepository, err := filesystem.NewFilesystemProblemCodeRepository(problemsDataDir)
	if err != nil {
		log.Fatal("failed to init problemCodeRepository", err)
	}

	dbTransactor := postgres.NewPostgreSQLTransactor(db)

	// services
	userService := services.NewUserService(
		userRepository,
		userProfileRepository,
		submissionRepository,
		problemRepository,
		mfaRepository,
		passwordResetRepository,
		emailVerificationRepository,
		dbTransactor,

		mailMan,
	)

	problemService := services.NewProblemService(
		problemRepository,
		problemNoteRepository,
		runRepository,
		submissionRepository,
		problemHintRepository,
		problemCodeRepository,

		judge,
	)

	staticService := services.NewStaticService()

	courseService := services.NewCourseService(
		courseRepository,
		lessonRepository,
		courseProgressRepository,
	)

	stripeService := services.NewStripeService(
		userRepository,
		stripeSecret,
	)

	// handlers
	userHandler := handlers.NewUserHandler(userService, staticService)
	problemHandler := handlers.NewProblemsHandlers(problemService)
	statsHandler := handlers.NewStatsHandler(userService, problemService)
	staticHandler := handlers.NewStaticHandler(staticService)
	adminHandler := handlers.NewAdminHandler(courseService, staticService, userService, problemService)
	courseHandler := handlers.NewCourseHandler(courseService)
	subscriptionHandler := handlers.NewSubscriptionHandler(stripeService)

	// middlewares
	authMiddleware := func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, userService)
	}
	partialAuthMiddleware := func(next http.Handler) http.Handler {
		return middleware.PartialAuthMiddleware(next, userService)
	}

	problemMiddleware := func(next http.Handler) http.Handler {
		return middleware.ProblemMiddleware(next, problemService)
	}
	courseMiddleware := func(next http.Handler) http.Handler {
		return middleware.CourseMiddleware(next, courseService)
	}
	lessonMiddleware := func(next http.Handler) http.Handler {
		return middleware.LessonMiddleware(next, courseService)
	}
	userMiddleware := func(next http.Handler) http.Handler {
		return middleware.UserMiddleware(next, userService)
	}
	hintMiddleware := func(next http.Handler) http.Handler {
		return middleware.ProblemHintMiddleware(next, problemService)
	}

	return router.Service(
		userHandler,
		problemHandler,
		statsHandler,
		staticHandler,
		adminHandler,
		courseHandler,
		subscriptionHandler,

		authMiddleware,
		partialAuthMiddleware,
		problemMiddleware,
		courseMiddleware,
		lessonMiddleware,
		userMiddleware,
		hintMiddleware,
	)
}
