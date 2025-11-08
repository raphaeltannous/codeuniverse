package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
}

type userService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		userRepo: r,
		logger:   slog.Default().With("package", "postgres.UserRepository"),
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

	return id, nil
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
	return nil, nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		slog.Error("failed to get user info by username", "err", err)
		return nil, fmt.Errorf("failed to get user by username")
	}

	return user, nil
}
