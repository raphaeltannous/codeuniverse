package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"git.riyt.dev/codeuniverse/internal/mailer"
	"git.riyt.dev/codeuniverse/internal/mailer/templates"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrEmptyEmail         = errors.New("email address cannot be empty")
	ErrEmptyUsername      = errors.New("username cannot be empty")
	ErrWeakPasswordLength = errors.New("password should be greater than 8")
)

type UserService interface {
	Create(ctx context.Context, username, password, email string) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetById(ctx context.Context, uuidString string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	GetAllUsers(ctx context.Context, offset, limit int) ([]*models.User, error)

	SendPasswordResetEmail(ctx context.Context, email string) error
	ResetPasswordByToken(ctx context.Context, token, newPassword string) error

	SendMfaCodeVerificationEmail(ctx context.Context, email string) error
	VerifyMfaCode(ctx context.Context, token, code string) error

	SendEmailVerificationEmail(ctx context.Context, email string) error
	VerifyEmailByToken(ctx context.Context, token string) error
}

type userService struct {
	userRepo              repository.UserRepository
	mfaRepo               repository.MfaCodeRepository
	passwordResetRepo     repository.PasswordResetRepository
	emailVerificationRepo repository.EmailVerificationRepository

	logger  *slog.Logger
	mailMan mailer.Mailer
}

func NewUserService(
	userRepo repository.UserRepository,
	mfaRepo repository.MfaCodeRepository,
	passwordResetRepo repository.PasswordResetRepository,
	emailVerificationRepo repository.EmailVerificationRepository,

	mailMan mailer.Mailer,
) UserService {
	return &userService{
		userRepo:              userRepo,
		mfaRepo:               mfaRepo,
		passwordResetRepo:     passwordResetRepo,
		emailVerificationRepo: emailVerificationRepo,

		logger:  slog.Default().With("package", "postgres.UserRepository"),
		mailMan: mailMan,
	}
}

func (s *userService) Create(ctx context.Context, username, password, email string) (uuid.UUID, error) {
	if email == "" {
		return uuid.UUID{}, ErrEmptyEmail
	}
	if username == "" {
		return uuid.UUID{}, ErrEmptyUsername
	}
	if len(password) < 8 {
		return uuid.UUID{}, ErrWeakPasswordLength
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash password", "err", err)
		return uuid.UUID{}, fmt.Errorf("failed to hash password")
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Email:        email,
		Role:         "user",
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error("creating user repo error", "err", err)
		return uuid.UUID{}, fmt.Errorf("service error creating user")
	}

	return id, s.SendEmailVerificationEmail(ctx, email)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *userService) GetById(ctx context.Context, id string) (*models.User, error) {
	if err := uuid.Validate(id); err != nil {
		return nil, fmt.Errorf("provided id is not a valid uuid: %w", err)
	}

	newId, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid from id: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, newId)
	if err != nil {
		return nil, fmt.Errorf("service error getting user info: %w", err)
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	users, err := s.userRepo.GetUsers(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("service failed to get from userRepo: %w", err)
	}

	return users, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		slog.Error("failed to get user info by email", "err", err)
		return nil, fmt.Errorf("failed to get user by email")
	}
	return user, nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		slog.Error("failed to get user info by username", "err", err)
		return nil, fmt.Errorf("failed to get user by username")
	}

	return user, nil
}

func (s *userService) SendPasswordResetEmail(ctx context.Context, email string) error {
	user, err := s.GetByEmail(ctx, email)
	fmt.Println(user, err)
	if err != nil {
		return err
	}

	token, err := utils.GenerateToken(64)
	if err != nil {
		return err
	}

	s.passwordResetRepo.Save(
		ctx,
		user.ID,
		utils.HashToken(token),
		time.Now().UTC().Add(10*time.Minute),
	)

	resetPasswordTmplData := templates.NewResetPasswordTmplData(
		user.Username,
		fmt.Sprintf("http://localhost:8080/accounts/password/reset?token=%s", token),
		"10",
	)

	var htmlBody bytes.Buffer
	err = templates.ResetPasswordTmpl.Execute(&htmlBody, resetPasswordTmplData)
	if err != nil {
		return err
	}

	return s.mailMan.SendHTML(
		ctx,
		email,
		"Password Reset Request",
		htmlBody.String(),
	)
}

func (s *userService) ResetPasswordByToken(ctx context.Context, token, newPassword string) error {
	passwordReset, err := s.passwordResetRepo.GetByTokenHash(ctx, utils.HashToken(token))
	if err != nil {
		return err
	}

	if !time.Now().UTC().Before(passwordReset.ExpiresAt) {
		return errors.New("time is expired")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdatePassword(
		ctx,
		passwordReset.UserId,
		hashedPassword,
	)
	if err != nil {
		return err
	}

	return s.passwordResetRepo.Save(
		ctx,
		passwordReset.UserId,
		passwordReset.Hash,
		time.Now().UTC(),
	)
}

func (s *userService) SendMfaCodeVerificationEmail(ctx context.Context, email string) error {

	return nil
}

func (s *userService) VerifyMfaCode(ctx context.Context, token, code string) error {
	return nil
}

func (s *userService) SendEmailVerificationEmail(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	token, err := utils.GenerateToken(32)
	if err != nil {
		return err
	}

	err = s.emailVerificationRepo.Save(
		ctx,
		user.ID,
		utils.HashToken(token),
		time.Now().UTC().Add(10*time.Minute),
	)
	if err != nil {
		slog.Error("failed to save email verification token to repo", "err", err)
		return fmt.Errorf("failed to save email verification")
	}

	verifyEmailTmplData := templates.NewVerifyEmailTmplData(
		user.Username,
		user.Email,
		fmt.Sprintf("http://localhost:8080/accounts/signup/email-verification?token=%s", token),
		"10",
	)

	var htmlBody bytes.Buffer
	err = templates.VerifyEmailTmpl.Execute(&htmlBody, verifyEmailTmplData)
	if err != nil {
		return err
	}

	return s.mailMan.SendHTML(
		ctx,
		email,
		"Email Verification",
		htmlBody.String(),
	)
}

func (s *userService) VerifyEmailByToken(ctx context.Context, token string) error {
	emailVerification, err := s.emailVerificationRepo.GetByTokenHash(ctx, utils.HashToken(token))
	if err != nil {
		return err
	}

	if !time.Now().UTC().Before(emailVerification.ExpiresAt) {
		return errors.New("time is expired")
	}

	// Token is valid
	// update
	err = s.userRepo.UpdateVerify(
		ctx,
		emailVerification.UserId,
		true,
	)
	if err != nil {
		return err
	}

	return s.emailVerificationRepo.Save(
		ctx,
		emailVerification.UserId,
		emailVerification.Hash,
		time.Now().UTC(),
	)
}
