package service

import (
	"context"
	"fmt"

	"backend-go/domain"
	"backend-go/internal/dto"
	"backend-go/internal/repository"
)

type ProductService interface {
	GetProducts(ctx context.Context, req dto.ProductQueryRequest) ([]domain.Product, int64, error)
	GetProductDetail(ctx context.Context, productID string) (*domain.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetProducts(ctx context.Context, req dto.ProductQueryRequest) ([]domain.Product, int64, error) {
	products, total, err := s.repo.GetAllProducts(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products, total, nil
}

func (s *productService) GetProductDetail(ctx context.Context, productID string) (*domain.Product, error) {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product detail: %w", err)
	}
	return product, nil
}
