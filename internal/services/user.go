package services

import (
	"context"
	"errors"
	"fmt"

	"git.riyt.dev/codeuniverse/internal/models"
	"git.riyt.dev/codeuniverse/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, username, password, email string) (uuid.UUID, error)
	GetUserInfoById(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) ([]*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		userRepo: r,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, password, email string) (uuid.UUID, error) {
	if email == "" {
		return uuid.UUID{}, errors.New("email address cannot be empty")
	}
	if username == "" {
		return uuid.UUID{}, errors.New("username cannot be empty")
	}
	if len(password) < 8 {
		return uuid.UUID{}, errors.New("password should be greater than 8")
	}

	// TODO: hash password
	user := &models.User{
		Username:     username,
		PasswordHash: password,
		Email:        email,
		Role:         "user",
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("service error creating user: %w", err)
	}

	return id, nil
}

func (s *userService) GetUserInfoById(ctx context.Context, id string) (*models.User, error) {
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
