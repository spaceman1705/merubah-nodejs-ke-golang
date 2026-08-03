package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"backend-go/domain"
)

type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID string) ([]domain.Cart, error)
	AddToCart(ctx context.Context, cart *domain.Cart) error
	UpdateQuantity(ctx context.Context, cartID string, quantity int) error
	DeleteCartItem(ctx context.Context, cartID string) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID string) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.WithContext(ctx).
		Preload("Product.Images").
		Preload("Store").
		Where("user_id = ?", userID).
		Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (r *cartRepository) AddToCart(ctx context.Context, cart *domain.Cart) error {
	// Cek produk dengan toko yang sama sudah ada di keranjang user
	var existingCart domain.Cart
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ? AND store_id = ?", cart.UserID, cart.ProductID, cart.StoreID).
		First(&existingCart).Error

	if err == nil {
		// update quantity-nya ditambahkan dengan yang baru
		existingCart.Quantity += cart.Quantity
		return r.db.WithContext(ctx).Save(&existingCart).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Jika belum buat record baru
		return r.db.WithContext(ctx).Create(cart).Error
	}

	return err
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, cartID string, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Cart{}).
		Where("id = ?", cartID).
		Update("quantity", quantity).Error
}

func (r *cartRepository) DeleteCartItem(ctx context.Context, cartID string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", cartID).
		Delete(&domain.Cart{}).Error
}
