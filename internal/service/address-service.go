package service

import (
	"context"
	"fmt"

	"backend-go/domain"
	"backend-go/internal/dto"
	"backend-go/internal/repository"
)

type AddressService interface {
	GetUserAddresses(ctx context.Context, userID string) ([]domain.UserAddress, error)
	AddAddress(ctx context.Context, userID string, req dto.CreateAddressRequest) error
	SetPrimaryAddress(ctx context.Context, userID string, addressID string) error
	DeleteAddress(ctx context.Context, addressID string, userID string) error
}

type addressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) GetUserAddresses(ctx context.Context, userID string) ([]domain.UserAddress, error) {
	addresses, err := s.repo.GetAddressesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}
	return addresses, nil
}

func (s *addressService) AddAddress(ctx context.Context, userID string, req dto.CreateAddressRequest) error {
	address := domain.UserAddress{
		UserID:        userID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Address:       req.Address,
		ProvinceID:    req.ProvinceID,
		CityID:        req.CityID,
		PostalCode:    req.PostalCode,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		IsMainAddress: req.IsMainAddress,
	}

	if err := s.repo.CreateAddress(ctx, &address); err != nil {
		return fmt.Errorf("failed to create address: %w", err)
	}
	return nil
}

func (s *addressService) SetPrimaryAddress(ctx context.Context, userID string, addressID string) error {
	if err := s.repo.SetPrimaryAddress(ctx, userID, addressID); err != nil {
		return fmt.Errorf("failed to set primary address: %w", err)
	}
	return nil
}

func (s *addressService) DeleteAddress(ctx context.Context, addressID string, userID string) error {
	if err := s.repo.DeleteAddress(ctx, addressID, userID); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	return nil
}
