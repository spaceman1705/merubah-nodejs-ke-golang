package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"backend-go/domain"
)

type OrderRepository interface {
	GetOrdersByUserID(ctx context.Context, userID string, status *string) ([]domain.Order, error)
	GetOrderByID(ctx context.Context, orderID string, userID string) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID string, status *string) ([]domain.Order, error) {
	var orders []domain.Order
	query := r.db.WithContext(ctx).
		Preload("OrderItems.Product.Images").
		Preload("UserAddress.UserCity").
		Preload("Store").
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, orderID string, userID string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems.Product.Images").
		Preload("UserAddress.UserCity").
		Preload("Store").
		First(&order, "id = ? AND user_id = ?", orderID, userID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	result := r.db.WithContext(ctx).Model(&domain.Order{}).Where("id = ?", orderID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}
