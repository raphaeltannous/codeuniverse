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
	RegisterUser(ctx context.Context, username, password, email string) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetById(ctx context.Context, uuidString string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetProfile(ctx context.Context, user *models.User) (*models.UserProfile, error)

	UpdateUserPatch(ctx context.Context, user *models.User, userUpdatePatch map[string]any) error

	GetAllUsers(ctx context.Context, getParams *repository.GetUsersParams) ([]*models.User, int, error)

	GetUsersCount(ctx context.Context) (int, error)
	GetUsersRegisteredLastNDaysCount(ctx context.Context, days int) (int, error)
	GetAdminCount(ctx context.Context) (int, error)

	GetRecentRegisteredUsers(ctx context.Context, limit int) ([]*models.User, error)

	UpdateUserProfilePatch(ctx context.Context, user *models.User, userProfileUpdatePatch map[string]string) error

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
	userProfileRepo       repository.UserProfileRepository
	submissionRepo        repository.SubmissionRepository
	problemRepo           repository.ProblemRepository
	mfaRepo               repository.MfaCodeRepository
	passwordResetRepo     repository.PasswordResetRepository
	emailVerificationRepo repository.EmailVerificationRepository

	dbTransactor repository.Transactor

	logger  *slog.Logger
	mailMan mailer.Mailer
}

func NewUserService(
	userRepo repository.UserRepository,
	userProfileRepo repository.UserProfileRepository,
	submissionRepo repository.SubmissionRepository,
	problemRepo repository.ProblemRepository,
	mfaRepo repository.MfaCodeRepository,
	passwordResetRepo repository.PasswordResetRepository,
	emailVerificationRepo repository.EmailVerificationRepository,

	dbTransactor repository.Transactor,

	mailMan mailer.Mailer,
) UserService {
	return &userService{
		userRepo:              userRepo,
		userProfileRepo:       userProfileRepo,
		submissionRepo:        submissionRepo,
		problemRepo:           problemRepo,
		mfaRepo:               mfaRepo,
		passwordResetRepo:     passwordResetRepo,
		emailVerificationRepo: emailVerificationRepo,

		dbTransactor: dbTransactor,

		logger:  slog.Default().With("package", "userService"),
		mailMan: mailMan,
	}
}

func (s *userService) RegisterUser(ctx context.Context, username, password, email string) (*models.User, error) {
	if !s.isEmailValid(email) {
		return nil, ErrInvalidEmail
	}

	if !s.isUsernameValid(username) {
		return nil, ErrInvalidUsername
	}

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

	err = s.dbTransactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		user, err = s.userRepo.Create(txCtx, user)
		if err != nil {
			return err
		}

		return s.userProfileRepo.Create(txCtx, user.ID)
	})

	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, err
		}

		s.logger.Error("creating user repo error", "err", err, "user.ID", user.ID)
		return nil, fmt.Errorf("service error creating user")
	}

	return user, s.SendEmailVerificationEmail(ctx, email)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *userService) UpdateUserPatch(ctx context.Context, user *models.User, userUpdatePatch map[string]any) error {
	if rawAvatarUrl, ok := userUpdatePatch["avatarUrl"]; ok {
		switch avatarUrl := rawAvatarUrl.(type) {
		case string:
			err := s.userRepo.UpdateAvatarUrl(ctx, user.ID, avatarUrl)
			if err != nil {
				s.logger.Error("failed to update avatarUrl", "user", user, "newAvatarUrl", avatarUrl, "err", err)
				return err
			}
		default:
			s.logger.Error("avatarUrl is not a string", "rawAvatarUrl", rawAvatarUrl)
			return ErrInvalidPatch
		}
	}

	if rawUsername, ok := userUpdatePatch["username"]; ok {
		switch username := rawUsername.(type) {
		case string:
			err := s.userRepo.UpdateUsername(ctx, user.ID, username)
			if err != nil {
				s.logger.Error("failed to update username", "user", user, "newUsername", username, "err", err)
				return err
			}
		default:
			s.logger.Error("username is not a string", "rawUsername", rawUsername)
			return ErrInvalidPatch
		}
	}

	if rawEmail, ok := userUpdatePatch["email"]; ok {
		switch email := rawEmail.(type) {
		case string:
			err := s.userRepo.UpdateEmail(ctx, user.ID, email)
			if err != nil {
				s.logger.Error("failed to update email", "user", user, "newEmail", email, "err", err)
				return err
			}
		default:
			s.logger.Error("email is not a string", "rawEmail", rawEmail)
			return ErrInvalidPatch
		}
	}

	if rawRole, ok := userUpdatePatch["role"]; ok {
		switch role := rawRole.(type) {
		case string:
			err := s.userRepo.UpdateRole(ctx, user.ID, role)
			if err != nil {
				s.logger.Error("failed to update role", "user", user, "newRole", role, "err", err)
				return err
			}
		default:
			s.logger.Error("role is not a string", "rawRole", rawRole)
			return ErrInvalidPatch
		}
	}

	if rawIsActive, ok := userUpdatePatch["isActive"]; ok {
		switch isActive := rawIsActive.(type) {
		case bool:
			err := s.userRepo.UpdateActive(ctx, user.ID, isActive)
			if err != nil {
				s.logger.Error("failed to update isActive", "user", user, "newIsActive", isActive, "err", err)
				return err
			}
		default:
			s.logger.Error("isActive is not a bool", "rawIsActive", rawIsActive)
			return ErrInvalidPatch
		}
	}

	if rawIsVerified, ok := userUpdatePatch["isVerified"]; ok {
		switch isVerified := rawIsVerified.(type) {
		case bool:
			err := s.userRepo.UpdateActive(ctx, user.ID, isVerified)
			if err != nil {
				s.logger.Error("failed to update isVerified", "user", user, "newIsVerified", isVerified, "err", err)
				return err
			}
		default:
			s.logger.Error("isVerified is not a bool", "rawIsVerified", rawIsVerified)
			return ErrInvalidPatch
		}
	}

	return nil
}

func (s *userService) GetAllUsers(ctx context.Context, getParams *repository.GetUsersParams) ([]*models.User, int, error) {
	users, total, err := s.userRepo.GetUsers(ctx, getParams)
	if err != nil {
		s.logger.Error("failed to get all users", "getParams", getParams, "err", err)
		return nil, 0, err
	}

	return users, total, nil
}

func (s *userService) GetAdminCount(ctx context.Context) (int, error) {
	count, err := s.userRepo.GetAdminCount(ctx)
	if err != nil {
		s.logger.Error("failed to get admin count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *userService) GetUsersCount(ctx context.Context) (int, error) {
	count, err := s.userRepo.GetUsersCount(ctx)
	if err != nil {
		s.logger.Error("failed to get users count", "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *userService) GetUsersRegisteredLastNDaysCount(ctx context.Context, since int) (int, error) {
	count, err := s.userRepo.GetUsersRegisteredLastNDaysCount(ctx, since)
	if err != nil {
		s.logger.Error("failed to get new users count", "since", since, "err", err)
		return 0, repository.ErrInternalServerError
	}

	return count, nil
}

func (s *userService) GetRecentRegisteredUsers(ctx context.Context, limit int) ([]*models.User, error) {
	users, err := s.userRepo.GetRecentRegisteredUsers(ctx, limit)
	if err != nil {
		s.logger.Error("failed to get recent registered users", "err", err)
		return nil, repository.ErrInternalServerError
	}

	return users, nil
}

func (s *userService) UpdateUserProfilePatch(ctx context.Context, user *models.User, userProfileUpdatePatch map[string]string) error {
	if avatarUrl, ok := userProfileUpdatePatch["avatarUrl"]; ok {
		if err := s.userRepo.UpdateAvatarUrl(ctx, user.ID, avatarUrl); err != nil {
			slog.Error("failed to update avatarUrl", "avatarUrl", avatarUrl, "user", user, "err", err)
			return err
		}
		return nil
	}

	if name, ok := userProfileUpdatePatch["name"]; ok {
		if err := s.userProfileRepo.UpdateName(ctx, user.ID, name); err != nil {
			slog.Error("failed to update name", "name", name, "user", user, "err", err)
			return err
		}
	}

	if bio, ok := userProfileUpdatePatch["bio"]; ok {
		if err := s.userProfileRepo.UpdateBio(ctx, user.ID, bio); err != nil {
			slog.Error("failed to update bio", "bio", bio, "user", user, "err", err)
			return err
		}
	}

	if preferredLanguage, ok := userProfileUpdatePatch["preferredLanguage"]; ok {
		if err := s.userProfileRepo.UpdatePreferredLanguage(ctx, user.ID, preferredLanguage); err != nil {
			slog.Error("failed to update preferredLanguage", "preferredLanguage", preferredLanguage, "user", user, "err", err)
			return err
		}
	}

	if country, ok := userProfileUpdatePatch["country"]; ok {
		if err := s.userProfileRepo.UpdateCountry(ctx, user.ID, country); err != nil {
			slog.Error("failed to update country", "country", country, "user", user, "err", err)
			return err
		}
	}

	if websiteUrl, ok := userProfileUpdatePatch["websiteUrl"]; ok {
		if err := s.userProfileRepo.UpdateWebsiteURL(ctx, user.ID, websiteUrl); err != nil {
			slog.Error("failed to update websiteUrl", "websiteUrl", websiteUrl, "user", user, "err", err)
			return err
		}
	}

	if linkedinUrl, ok := userProfileUpdatePatch["linkedinUrl"]; ok {
		if err := s.userProfileRepo.UpdateLinkedinURL(ctx, user.ID, linkedinUrl); err != nil {
			slog.Error("failed to update linkedinUrl", "linkedinUrl", linkedinUrl, "user", user, "err", err)
			return err
		}
	}

	if xUrl, ok := userProfileUpdatePatch["xUrl"]; ok {
		if err := s.userProfileRepo.UpdateXURL(ctx, user.ID, xUrl); err != nil {
			slog.Error("failed to update xUrl", "xUrl", xUrl, "user", user, "err", err)
			return err
		}
	}

	if githubUrl, ok := userProfileUpdatePatch["githubUrl"]; ok {
		if err := s.userProfileRepo.UpdateGithubURL(ctx, user.ID, githubUrl); err != nil {
			slog.Error("failed to update githubUrl", "githubUrl", githubUrl, "user", user, "err", err)
			return err
		}
	}

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

func (s *userService) GetProfile(ctx context.Context, user *models.User) (*models.UserProfile, error) {
	userProfile, err := s.userProfileRepo.GetInfo(ctx, user)
	if err != nil {
		s.logger.Error("failed to get userProfile", "err", err, "user", user)
		return nil, err
	}

	userProfile.AvatarURL = user.AvatarURL

	submissionStats, err := s.submissionRepo.GetSubmissionsStats(ctx, user.ID)
	if err != nil {
		s.logger.Error("failed to get submissionStats", "err", err)
		return nil, err
	}
	userProfile.SubmissionStats = *submissionStats

	userProfile.EasyCount, err = s.problemRepo.GetEasyCount(ctx)
	if err != nil {
		s.logger.Error("failed to get easy count", "err", err)
		return nil, err
	}

	userProfile.MediumCount, err = s.problemRepo.GetMediumCount(ctx)
	if err != nil {
		s.logger.Error("failed to get medium count", "err", err)
		return nil, err
	}

	userProfile.HardCount, err = s.problemRepo.GetHardCount(ctx)
	if err != nil {
		s.logger.Error("failed to get hard count", "err", err)
		return nil, err
	}

	return userProfile, err
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

	err = s.passwordResetRepo.Save(
		ctx,
		user.ID,
		utils.HashToken(token),
		time.Now().UTC().Add(10*time.Minute),
	)
	if err != nil {
		return err
	}

	resetPasswordTmplData := templates.NewResetPasswordTmplData(
		user.Username,
		// TODO: is there a better way to point to application url?
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

	newToken, err := utils.GenerateToken(16)
	if err != nil {
		return err
	}

	return s.passwordResetRepo.Save(
		ctx,
		passwordReset.UserId,
		utils.HashToken(newToken),
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
	// TODO: I think these functions should return map[string]string or an error.
	// for example isPasswordValid:
	// -> errors.New("weak password lenght")
	// -> errors.New("weak password")
	// -> errors.New("password does not validate all rules")
	// and so on.
	return false
}
