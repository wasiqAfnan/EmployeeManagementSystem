package service

import (
	"context"
	"errors"

	"EMS/internal/model"
	"EMS/internal/repository"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("user with this email already exists")
	ErrDuplicateSub   = errors.New("user with this sub already exists")
	ErrAdminRequired  = errors.New("admin access required")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserBySub(ctx context.Context, sub string) (*model.User, error) {
	user, err := s.repo.FindBySub(ctx, sub)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, role string) ([]model.User, error) {
	if role != "admin" {
		return nil, ErrAdminRequired
	}
	return s.repo.GetAll(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	// Check duplicate sub
	existingSub, err := s.repo.FindBySub(ctx, user.Sub)
	if err != nil {
		return nil, err
	}
	if existingSub != nil {
		return nil, ErrDuplicateSub
	}

	// Check duplicate email
	existingEmail, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, ErrDuplicateEmail
	}

	return s.repo.Create(ctx, user)
}
