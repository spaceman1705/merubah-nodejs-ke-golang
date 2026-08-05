package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"backend-go/domain"
)

type AddressRepository interface {
	GetAddressesByUserID(ctx context.Context, userID string) ([]domain.UserAddress, error)
	CreateAddress(ctx context.Context, address *domain.UserAddress) error
	SetPrimaryAddress(ctx context.Context, userID string, addressID string) error
	DeleteAddress(ctx context.Context, addressID string, userID string) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) GetAddressesByUserID(ctx context.Context, userID string) ([]domain.UserAddress, error) {
	var addresses []domain.UserAddress
	err := r.db.WithContext(ctx).
		Preload("Province").
		Preload("UserCity").
		Where("user_id = ?", userID).
		Order("is_main_address desc, created_at desc").
		Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) CreateAddress(ctx context.Context, address *domain.UserAddress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if address.IsMainAddress {
			if err := tx.Model(&domain.UserAddress{}).Where("user_id = ?", address.UserID).Update("is_main_address", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(address).Error
	})
}

func (r *addressRepository) SetPrimaryAddress(ctx context.Context, userID string, addressID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.UserAddress{}).Where("user_id = ?", userID).Update("is_main_address", false).Error; err != nil {
			return err
		}
		result := tx.Model(&domain.UserAddress{}).Where("id = ? AND user_id = ?", addressID, userID).Update("is_main_address", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("address not found or unauthorized")
		}
		return nil
	})
}

func (r *addressRepository) DeleteAddress(ctx context.Context, addressID string, userID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", addressID, userID).Delete(&domain.UserAddress{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("address not found or unauthorized")
	}
	return nil
}
