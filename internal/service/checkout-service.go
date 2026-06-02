package service

import (
	"context"
	"errors"
	"fmt"

	"backend-go/domain"
	"backend-go/internal/dto"
	"backend-go/internal/repository"
)

type CheckoutService interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*domain.Order, error)
	GetUserOrders(ctx context.Context, userID string) ([]domain.Order, error)
	CancelOrder(ctx context.Context, orderID string, userID string, reason string) (*domain.Order, error)
}

type checkoutService struct {
	repo repository.CheckoutRepository
}

func NewCheckoutService(repo repository.CheckoutRepository) CheckoutService {
	return &checkoutService{repo: repo}
}

func (s *checkoutService) GetUserOrders(ctx context.Context, userID string) ([]domain.Order, error) {
	orders, err := s.repo.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	return orders, nil
}

func (s *checkoutService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*domain.Order, error) {
	_, _, err := s.repo.CheckAddressAndStore(ctx, req.UserID, req.UserAddressID, req.StoreID)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	computedSubtotal := 0.0
	for _, item := range req.Items {
		product, stock, err := s.repo.GetProductWithStock(ctx, item.ProductID, req.StoreID)
		if err != nil {
			return nil, err
		}
		if product == nil || !product.IsActive {
			return nil, fmt.Errorf("product %s is not available", item.ProductID)
		}
		if stock == nil || stock.Quantity < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}
		computedSubtotal += (item.Price * float64(item.Quantity))
	}

	computedDiscount := req.DiscountAmount
	itemDiscounts := make(map[string]float64)
	usageMap := make(map[string]float64)

	expectedTotal := computedSubtotal + req.ShippingCost - computedDiscount
	if expectedTotal < 0 {
		expectedTotal = 0
	}

	if req.TotalAmount != expectedTotal {
		return nil, errors.New("total amount mismatch, possible data manipulation")
	}

	order, err := s.repo.CreateOrderTransaction(ctx, req, computedSubtotal, computedDiscount, itemDiscounts, usageMap)
	if err != nil {
		return nil, fmt.Errorf("failed to process checkout transaction: %w", err)
	}

	return order, nil
}

func (s *checkoutService) CancelOrder(ctx context.Context, orderID string, userID string, reason string) (*domain.Order, error) {
	if reason == "" {
		reason = "Dibatalkan oleh pengguna (pembeli)"
	}

	order, err := s.repo.CancelOrder(ctx, orderID, userID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	return order, nil
}
