package repository

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// ProductRepository defines the interface for product persistence operations.
// Following the Interface Segregation Principle (ISP) from SOLID.
type ProductRepository interface {
	// Create stores a new product.
	Create(ctx context.Context, product *entity.Product) error
	// GetByID retrieves a product by its ID.
	GetByID(ctx context.Context, id string) (*entity.Product, error)
	// GetAll retrieves all products with pagination.
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	// GetByCategory retrieves products by category ID.
	GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Product, error)
	// Update updates an existing product.
	Update(ctx context.Context, product *entity.Product) error
	// Delete removes a product by its ID.
	Delete(ctx context.Context, id string) error
	// Search searches products by name or description.
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error)
	// GetActiveProducts retrieves all active products.
	GetActiveProducts(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	// UpdateStock updates the stock of a product.
	UpdateStock(ctx context.Context, id string, stock int) error
	// Count returns the total number of products.
	Count(ctx context.Context) (int64, error)
}
