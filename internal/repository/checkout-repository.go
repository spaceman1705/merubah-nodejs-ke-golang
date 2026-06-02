package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend-go/domain"
	"backend-go/internal/dto"
)

type CheckoutRepository interface {
	CheckAddressAndStore(ctx context.Context, userID, addressID, storeID string) (*domain.UserAddress, *domain.Store, error)
	GetProductWithStock(ctx context.Context, productID, storeID string) (*domain.Product, *domain.ProductStock, error)
	CreateOrderTransaction(ctx context.Context, req *dto.CreateOrderRequest, computedSubtotal, computedDiscount float64, itemDiscounts map[string]float64, usageMap map[string]float64) (*domain.Order, error)
	GetUserOrders(ctx context.Context, userID string) ([]domain.Order, error)
	CancelOrder(ctx context.Context, orderID string, userID string, cancelReason string) (*domain.Order, error)
}

type checkoutRepository struct {
	db *gorm.DB
}

func NewCheckoutRepository(db *gorm.DB) CheckoutRepository {
	return &checkoutRepository{db: db}
}

func (r *checkoutRepository) CheckAddressAndStore(ctx context.Context, userID, addressID, storeID string) (*domain.UserAddress, *domain.Store, error) {
	var address domain.UserAddress
	if err := r.db.WithContext(ctx).First(&address, "id = ? AND user_id = ?", addressID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("address not found or does not belong to user")
		}
		return nil, nil, err
	}

	var store domain.Store
	if err := r.db.WithContext(ctx).First(&store, "id = ?", storeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("store not found")
		}
		return nil, nil, err
	}

	return &address, &store, nil
}

func (r *checkoutRepository) GetProductWithStock(ctx context.Context, productID, storeID string) (*domain.Product, *domain.ProductStock, error) {
	var product domain.Product
	if err := r.db.WithContext(ctx).Preload("Stocks", "store_id = ?", storeID).First(&product, "id = ?", productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("product %s not found", productID)
		}
		return nil, nil, err
	}

	if len(product.Stocks) == 0 {
		return &product, nil, nil
	}

	return &product, &product.Stocks[0], nil
}

func (r *checkoutRepository) GetUserOrders(ctx context.Context, userID string) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.WithContext(ctx).
		Preload("OrderItems.Product.Images").
		Preload("UserAddress.UserCity").
		Preload("Store").
		Order("created_at desc").
		Find(&orders, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *checkoutRepository) CreateOrderTransaction(
	ctx context.Context,
	req *dto.CreateOrderRequest,
	computedSubtotal float64,
	computedDiscount float64,
	itemDiscounts map[string]float64,
	usageMap map[string]float64,
) (*domain.Order, error) {

	var newOrder domain.Order

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderNumber := generateOrderNumber()

		initialStatus := "WAITING_CONFIRMATION" // Gunakan string langsung untuk amannya jika enum berbeda
		if req.PaymentMethod == "TRANSFER" {
			initialStatus = "WAITING_PAYMENT"
		}

		newOrder = domain.Order{
			OrderNumber:      orderNumber,
			UserID:           req.UserID,
			UserAddressID:    req.UserAddressID,
			StoreID:          req.StoreID,
			Subtotal:         computedSubtotal,
			ShippingCost:     req.ShippingCost,
			DiscountAmount:   computedDiscount,
			ShippingDiscount: 0,
			TotalAmount:      req.TotalAmount,
			VoucherCodeUsed:  req.VoucherCodeUsed,
			Status:           domain.OrderStatus(initialStatus),
		}

		if err := tx.Create(&newOrder).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		for _, item := range req.Items {
			lineSubtotal := (item.Price * float64(item.Quantity))
			itemDiscount := itemDiscounts[item.ProductID]

			orderItem := domain.OrderItem{
				OrderID:        newOrder.ID,
				ProductID:      item.ProductID,
				Quantity:       item.Quantity,
				Price:          item.Price,
				DiscountAmount: itemDiscount,
				Subtotal:       lineSubtotal - itemDiscount,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				return fmt.Errorf("failed to create order item for product %s: %w", item.ProductID, err)
			}

			var productStock domain.ProductStock
			if err := tx.First(&productStock, "product_id = ? AND store_id = ?", item.ProductID, req.StoreID).Error; err != nil {
				return fmt.Errorf("failed to find product stock: %w", err)
			}

			if err := tx.Model(&productStock).UpdateColumn("quantity", gorm.Expr("quantity - ?", item.Quantity)).Error; err != nil {
				return fmt.Errorf("failed to decrement stock: %w", err)
			}

			journal := domain.StockJournal{
				ProductStockID: productStock.ID,
				Action:         "OUT",
				Quantity:       item.Quantity,
			}
			if err := tx.Create(&journal).Error; err != nil {
				return fmt.Errorf("failed to create stock journal: %w", err)
			}
		}

		if err := tx.Where("user_id = ?", req.UserID).Delete(&domain.Cart{}).Error; err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var orderWithDetails domain.Order
	err = r.db.WithContext(ctx).
		Preload("OrderItems.Product.Images").
		Preload("UserAddress.UserCity").
		Preload("Store").
		First(&orderWithDetails, "id = ?", newOrder.ID).Error

	return &orderWithDetails, err
}

func generateOrderNumber() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	uuidStr := uuid.NewString()[:8]
	return fmt.Sprintf("ORD-%d-%s", timestamp, uuidStr)
}

func (r *checkoutRepository) CancelOrder(ctx context.Context, orderID string, userID string, cancelReason string) (*domain.Order, error) {
	var order domain.Order

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("OrderItems").First(&order, "id = ? AND user_id = ?", orderID, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("order not found or unauthorized")
			}
			return err
		}

		if string(order.Status) == "CANCELLED" || string(order.Status) == "SHIPPED" || string(order.Status) == "COMPLETED" {
			return errors.New("order cannot be cancelled because of its current status")
		}

		now := time.Now()
		order.Status = domain.OrderStatus("CANCELLED")
		order.CancelledAt = &now
		order.CancellationReason = &cancelReason

		if err := tx.Save(&order).Error; err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		for _, item := range order.OrderItems {
			var productStock domain.ProductStock
			if err := tx.First(&productStock, "product_id = ? AND store_id = ?", item.ProductID, order.StoreID).Error; err != nil {
				return fmt.Errorf("failed to find product stock to restore: %w", err)
			}

			if err := tx.Model(&productStock).UpdateColumn("quantity", gorm.Expr("quantity + ?", item.Quantity)).Error; err != nil {
				return fmt.Errorf("failed to restore stock: %w", err)
			}

			journal := domain.StockJournal{
				ProductStockID: productStock.ID,
				Action:         "IN",
				Quantity:       item.Quantity,
			}
			if err := tx.Create(&journal).Error; err != nil {
				return fmt.Errorf("failed to create stock journal for restoration: %w", err)
			}
		}

		return nil
	})

	return &order, err
}
