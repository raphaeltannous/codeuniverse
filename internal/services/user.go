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
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("service error creating user: %w", err)
	}

	return id, nil
}
