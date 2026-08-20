package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"backend-go/domain"
)

type UserRepository interface {
	GetProfileByID(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetProfileByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Omit("Password"). // tidak ambil password
		First(&user, "id = ?", userID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) UpdateProfile(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Model(user).Updates(user).Error
}
