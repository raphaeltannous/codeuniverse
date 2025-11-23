package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"time"

	"git.riyt.dev/codeuniverse/internal/mailer"
	"git.riyt.dev/codeuniverse/internal/mailer/templates"
	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"git.riyt.dev/codeuniverse/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrWeakPasswordLength = errors.New("password should be greater than 8")

	ErrTimeIsExpired = errors.New("time is expired")

	ErrInvalidMfaCode = errors.New("invalid mfa code")
)

type UserService interface {
	Create(ctx context.Context, username, password, email string) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetById(ctx context.Context, uuidString string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	GetAllUsers(ctx context.Context, offset, limit int) ([]*models.User, error)

	SendPasswordResetEmail(ctx context.Context, email string) error
	ResetPasswordByToken(ctx context.Context, token, newPassword string) error

	GetMfaCodeByToken(ctx context.Context, token string) (*models.MfaCode, error)
	CreateMfaCodeAndToken(ctx context.Context, user *models.User) (string, string, error)
	SendMfaCodeVerificationEmail(ctx context.Context, user *models.User, mfaCode string) error
	VerifyMfaCode(ctx context.Context, token, code string) (*models.MfaCode, error)

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

func (s *userService) Create(ctx context.Context, username, password, email string) (*models.User, error) {
	if !s.isEmailValid(email) {
		return nil, ErrInvalidEmail
	}

	// TODO: Validate username
	// TODO: Validate password

	// if email == "" {
	// 	return uuid.UUID{}, ErrEmptyEmail
	// }

	// if username == "" {
	// 	return uuid.UUID{}, ErrEmptyUsername
	// }

	if len(password) < 8 {
		return nil, ErrWeakPasswordLength
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash password", "err", err)
		return nil, fmt.Errorf("failed to hash password")
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Email:        email,
		Role:         "user",
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, err
		}

		s.logger.Error("creating user repo error", "err", err)
		return nil, fmt.Errorf("service error creating user")
	}

	return user, s.SendEmailVerificationEmail(ctx, email)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *userService) GetAllUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	users, err := s.userRepo.GetUsers(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("service failed to get from userRepo: %w", err)
	}

	return users, nil
}

func (s *userService) GetById(ctx context.Context, id string) (*models.User, error) {
	if err := uuid.Validate(id); err != nil {
		return nil, fmt.Errorf("provided id is not a valid uuid: %w", err)
	}

	newId, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid from id: %w", err)
	}

	getFn := func(ctx context.Context) (*models.User, error) {
		return s.userRepo.GetByID(ctx, newId)
	}

	return s.getByFunc(ctx, getFn)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	getFn := func(ctx context.Context) (*models.User, error) {
		return s.userRepo.GetByEmail(ctx, email)
	}

	return s.getByFunc(ctx, getFn)
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	getFn := func(ctx context.Context) (*models.User, error) {
		return s.userRepo.GetByUsername(ctx, username)
	}

	return s.getByFunc(ctx, getFn)
}

func (s *userService) getByFunc(ctx context.Context, getFn func(ctx context.Context) (*models.User, error)) (*models.User, error) {
	user, err := getFn(ctx)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, err
		default:
			s.logger.Error("failed to get user", "err", err, "fn", getFn)

			return nil, fmt.Errorf("internal server error")
		}
	}

	return user, err
}

func (s *userService) SendPasswordResetEmail(ctx context.Context, email string) error {
	user, err := s.GetByEmail(ctx, email)
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
		return ErrTimeIsExpired
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

func (s *userService) GetMfaCodeByToken(ctx context.Context, token string) (*models.MfaCode, error) {
	mfaCode, err := s.mfaRepo.GetByTokenHash(
		ctx,
		utils.HashToken(token),
	)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrMfaTokenNotFound):
			return nil, err
		default:
			s.logger.Error("GetMfaCodeByToken failed to get mfaCode", "err", err)
			return nil, repository.ErrInternalServerError
		}
	}

	return mfaCode, nil
}

func (s *userService) CreateMfaCodeAndToken(ctx context.Context, user *models.User) (string, string, error) {
	token, err := utils.GenerateToken(32)
	if err != nil {
		return "", "", err
	}

	mfaCode, err := utils.GenerateNumericCode(7)
	if err != nil {
		return "", "", err
	}

	err = s.mfaRepo.Save(
		ctx,
		user.ID,
		utils.HashToken(token),
		utils.HashToken(mfaCode),
		time.Now().UTC().Add(10*time.Minute),
	)

	if err != nil {
		slog.Error("failed to save mfa code to repo", "err", err)
		return "", "", fmt.Errorf("failed to save mfa code to repo")
	}

	return mfaCode, token, nil
}

func (s *userService) SendMfaCodeVerificationEmail(ctx context.Context, user *models.User, mfaCode string) error {
	mfaTmplData := templates.NewTwoFATmplData(
		user.Username,
		mfaCode,
		"10",
	)

	var htmlBody bytes.Buffer
	err := templates.TwoFATmpl.Execute(&htmlBody, mfaTmplData)
	if err != nil {
		return err
	}

	return s.mailMan.SendHTML(
		ctx,
		user.Email,
		"MFA Verification",
		htmlBody.String(),
	)
}

func (s *userService) VerifyMfaCode(ctx context.Context, token, code string) (*models.MfaCode, error) {
	mfaCode, err := s.mfaRepo.GetByTokenHash(
		ctx,
		utils.HashToken(token),
	)

	if err != nil {
		if errors.Is(err, repository.ErrMfaTokenNotFound) {
			return nil, err
		}

		return nil, repository.ErrInternalServerError
	}

	if !time.Now().UTC().Before(mfaCode.ExpiresAt) {
		return nil, ErrTimeIsExpired
	}

	if codeHash := utils.HashToken(code); codeHash != mfaCode.CodeHash {
		s.logger.Debug("invalid code hash", "codeHash", codeHash, "mfaCode.CodeHash", mfaCode.CodeHash)
		return nil, ErrInvalidMfaCode
	}

	newToken, err := utils.GenerateToken(32)
	if err != nil {
		return nil, err
	}

	err = s.mfaRepo.Save(
		ctx,
		mfaCode.UserId,
		utils.HashToken(newToken),
		mfaCode.CodeHash,
		time.Now().UTC(),
	)

	return mfaCode, err
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
		s.logger.Debug("emailVerification", "emailVerification", emailVerification, "err", err)
		return err
	}

	if !time.Now().UTC().Before(emailVerification.ExpiresAt) {
		return ErrTimeIsExpired
	}

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

func (s *userService) isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *userService) isUsernameValid(username string) bool {
	if len(username) < 3 || len(username) > 25 {
		return false
	}

	return true
}

func (s *userService) isPasswordValid(password string) bool {
	return false
}
