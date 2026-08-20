package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"backend-go/domain"
	"backend-go/internal/dto"
	"backend-go/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*domain.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	// Cek apakah email sudah terdaftar
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return errors.New("email is already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	passStr := string(hashedPassword)

	// Buat referral code sederhana
	referralCode := generateReferralCode()

	newUser := domain.User{
		Email:        req.Email,
		Password:     &passStr,
		ReferralCode: referralCode,
		IsVerified:   false, // Set true jika tidak pakai verifikasi email
		IsActive:     true,
		Role:         domain.RoleUser,
	}

	if err := s.userRepo.CreateUser(ctx, &newUser); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Password == nil {
		return nil, errors.New("invalid credentials")
	}

	// Bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Sembunyikan password di kembalian
	user.Password = nil
	return user, nil
}

func generateReferralCode() string {
	rand.Seed(time.Now().UnixNano())
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
