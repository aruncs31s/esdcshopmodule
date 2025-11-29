package persistence

import (
	"context"
	"strings"
	"sync"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
)

// InMemoryProductRepository is an in-memory implementation of ProductRepository.
// This is useful for testing and development.
type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[string]*entity.Product
}

// NewInMemoryProductRepository creates a new InMemoryProductRepository.
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*entity.Product),
	}
}

// Ensure InMemoryProductRepository implements ProductRepository.
var _ repository.ProductRepository = (*InMemoryProductRepository)(nil)

// Create stores a new product.
func (r *InMemoryProductRepository) Create(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; exists {
		return entity.ErrInvalidProductID
	}

	r.products[product.ID] = product
	return nil
}

// GetByID retrieves a product by its ID.
func (r *InMemoryProductRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, entity.ErrProductNotFound
	}
	return product, nil
}

// GetAll retrieves all products with pagination.
func (r *InMemoryProductRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*entity.Product, 0, len(r.products))
	for _, p := range r.products {
		products = append(products, p)
	}

	// Apply pagination
	if offset >= len(products) {
		return []*entity.Product{}, nil
	}

	end := offset + limit
	if end > len(products) {
		end = len(products)
	}

	return products[offset:end], nil
}

// GetByCategory retrieves products by category ID.
func (r *InMemoryProductRepository) GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*entity.Product, 0)
	for _, p := range r.products {
		if p.CategoryID == categoryID {
			products = append(products, p)
		}
	}

	// Apply pagination
	if offset >= len(products) {
		return []*entity.Product{}, nil
	}

	end := offset + limit
	if end > len(products) {
		end = len(products)
	}

	return products[offset:end], nil
}

// Update updates an existing product.
func (r *InMemoryProductRepository) Update(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return entity.ErrProductNotFound
	}

	r.products[product.ID] = product
	return nil
}

// Delete removes a product by its ID.
func (r *InMemoryProductRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return entity.ErrProductNotFound
	}

	delete(r.products, id)
	return nil
}

// Search searches products by name or description.
func (r *InMemoryProductRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(query)
	products := make([]*entity.Product, 0)

	for _, p := range r.products {
		if strings.Contains(strings.ToLower(p.Name), query) ||
			strings.Contains(strings.ToLower(p.Description), query) {
			products = append(products, p)
		}
	}

	// Apply pagination
	if offset >= len(products) {
		return []*entity.Product{}, nil
	}

	end := offset + limit
	if end > len(products) {
		end = len(products)
	}

	return products[offset:end], nil
}

// GetActiveProducts retrieves all active products.
func (r *InMemoryProductRepository) GetActiveProducts(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*entity.Product, 0)
	for _, p := range r.products {
		if p.Active {
			products = append(products, p)
		}
	}

	// Apply pagination
	if offset >= len(products) {
		return []*entity.Product{}, nil
	}

	end := offset + limit
	if end > len(products) {
		end = len(products)
	}

	return products[offset:end], nil
}

// UpdateStock updates the stock of a product.
func (r *InMemoryProductRepository) UpdateStock(ctx context.Context, id string, stock int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, exists := r.products[id]
	if !exists {
		return entity.ErrProductNotFound
	}

	return product.UpdateStock(stock)
}

// Count returns the total number of products.
func (r *InMemoryProductRepository) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return int64(len(r.products)), nil
}
