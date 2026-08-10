package service

import (
	"context"
	"fmt"

	"backend-go/domain"
	"backend-go/internal/repository"
)

type OrderService interface {
	GetUserOrders(ctx context.Context, userID string, status *string) ([]domain.Order, error)
	GetOrderDetail(ctx context.Context, orderID string, userID string) (*domain.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) GetUserOrders(ctx context.Context, userID string, status *string) ([]domain.Order, error) {
	orders, err := s.repo.GetOrdersByUserID(ctx, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user orders: %w", err)
	}
	return orders, nil
}

func (s *orderService) GetOrderDetail(ctx context.Context, orderID string, userID string) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order detail: %w", err)
	}
	return order, nil
}
