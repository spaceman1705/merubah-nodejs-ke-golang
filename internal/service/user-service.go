package service

import (
	"context"
	"fmt"

	"backend-go/domain"
	"backend-go/internal/dto"
	"backend-go/internal/repository"
)

type UserService interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*domain.User, error) {
	user, err := s.repo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.AvatarID != nil {
		user.AvatarID = req.AvatarID
	}

	if err := s.repo.UpdateProfile(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// hide password sebelum kembalikan ke response
	user.Password = nil
	return user, nil
}
