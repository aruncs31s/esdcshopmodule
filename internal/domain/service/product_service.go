package service

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
)

// ProductService provides domain-level operations for products.
type ProductService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new ProductService.
func NewProductService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// ValidateProductCategory validates that a product's category exists.
func (s *ProductService) ValidateProductCategory(ctx context.Context, categoryID string) error {
	if categoryID == "" {
		return nil // Category is optional
	}

	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return entity.ErrCategoryNotFound
	}
	return nil
}

// GetProductsWithCategory retrieves products with their category information.
func (s *ProductService) GetProductsWithCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Product, *entity.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, nil, err
	}

	products, err := s.productRepo.GetByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	return products, category, nil
}

// CheckStockAvailability checks if a product has enough stock.
func (s *ProductService) CheckStockAvailability(ctx context.Context, productID string, quantity int) (bool, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return false, err
	}
	return product.Stock >= quantity, nil
}

// GetAvailableProducts retrieves all products that are active and in stock.
func (s *ProductService) GetAvailableProducts(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	products, err := s.productRepo.GetActiveProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Filter products with stock > 0
	available := make([]*entity.Product, 0)
	for _, p := range products {
		if p.IsAvailable() {
			available = append(available, p)
		}
	}
	return available, nil
}
