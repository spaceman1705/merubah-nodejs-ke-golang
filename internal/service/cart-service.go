package service

import (
	"context"
	"errors"
	"fmt"

	"backend-go/domain"
	"backend-go/internal/dto"
	"backend-go/internal/repository"
)

type CartService interface {
	GetCartByUserID(ctx context.Context, userID string) ([]domain.Cart, error)
	AddToCart(ctx context.Context, userID string, req dto.AddToCartRequest) error
	UpdateCartQuantity(ctx context.Context, cartID string, userID string, quantity int) error
	RemoveCartItem(ctx context.Context, cartID string, userID string) error
}

type cartService struct {
	repo     repository.CartRepository
	checkout repository.CheckoutRepository //validasi stok produk
}

func NewCartService(repo repository.CartRepository, checkout repository.CheckoutRepository) CartService {
	return &cartService{repo: repo, checkout: checkout}
}

func (s *cartService) GetCartByUserID(ctx context.Context, userID string) ([]domain.Cart, error) {
	carts, err := s.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}
	return carts, nil
}

func (s *cartService) AddToCart(ctx context.Context, userID string, req dto.AddToCartRequest) error {
	// Validasi produk dan stoknya ada di toko
	_, stock, err := s.checkout.GetProductWithStock(ctx, req.ProductID, req.StoreID)
	if err != nil {
		return fmt.Errorf("product or stock check failed: %w", err)
	}
	if stock == nil || stock.Quantity < req.Quantity {
		return errors.New("insufficient stock for this product")
	}

	cartItem := &domain.Cart{
		UserID:    userID,
		ProductID: req.ProductID,
		StoreID:   req.StoreID,
		Quantity:  req.Quantity,
	}

	if err := s.repo.AddToCart(ctx, cartItem); err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}

	return nil
}

func (s *cartService) UpdateCartQuantity(ctx context.Context, cartID string, userID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if err := s.repo.UpdateQuantity(ctx, cartID, quantity); err != nil {
		return fmt.Errorf("failed to update cart quantity: %w", err)
	}

	return nil
}

func (s *cartService) RemoveCartItem(ctx context.Context, cartID string, userID string) error {
	if err := s.repo.DeleteCartItem(ctx, cartID); err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}
	return nil
}
