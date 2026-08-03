package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"backend-go/domain"
	"backend-go/internal/dto"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context, req dto.ProductQueryRequest) ([]domain.Product, int64, error)
	GetProductByID(ctx context.Context, productID string) (*domain.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAllProducts(ctx context.Context, req dto.ProductQueryRequest) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{}).Where("is_active = ?", true)

	// Filter pencarian
	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	// Filter
	if req.Category != "" {
		query = query.Where("product_category_id = ?", req.Category)
	}

	// Hitung total data keperluan pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.
		Preload("Images").
		Preload("Category").
		Preload("Stocks").
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, productID string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Images").
		Preload("Category").
		Preload("Stocks.Store").
		First(&product, "id = ?", productID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}
