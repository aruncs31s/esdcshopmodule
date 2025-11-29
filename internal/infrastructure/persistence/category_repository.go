package persistence

import (
	"context"
	"sync"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
)

// InMemoryCategoryRepository is an in-memory implementation of CategoryRepository.
type InMemoryCategoryRepository struct {
	mu         sync.RWMutex
	categories map[string]*entity.Category
}

// NewInMemoryCategoryRepository creates a new InMemoryCategoryRepository.
func NewInMemoryCategoryRepository() *InMemoryCategoryRepository {
	return &InMemoryCategoryRepository{
		categories: make(map[string]*entity.Category),
	}
}

// Ensure InMemoryCategoryRepository implements CategoryRepository.
var _ repository.CategoryRepository = (*InMemoryCategoryRepository)(nil)

// Create stores a new category.
func (r *InMemoryCategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[category.ID]; exists {
		return entity.ErrInvalidCategoryID
	}

	r.categories[category.ID] = category
	return nil
}

// GetByID retrieves a category by its ID.
func (r *InMemoryCategoryRepository) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	category, exists := r.categories[id]
	if !exists {
		return nil, entity.ErrCategoryNotFound
	}
	return category, nil
}

// GetAll retrieves all categories.
func (r *InMemoryCategoryRepository) GetAll(ctx context.Context) ([]*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]*entity.Category, 0, len(r.categories))
	for _, c := range r.categories {
		categories = append(categories, c)
	}
	return categories, nil
}

// GetByParentID retrieves categories by parent ID.
func (r *InMemoryCategoryRepository) GetByParentID(ctx context.Context, parentID string) ([]*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]*entity.Category, 0)
	for _, c := range r.categories {
		if c.ParentID == parentID {
			categories = append(categories, c)
		}
	}
	return categories, nil
}

// GetRootCategories retrieves all root categories (no parent).
func (r *InMemoryCategoryRepository) GetRootCategories(ctx context.Context) ([]*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]*entity.Category, 0)
	for _, c := range r.categories {
		if c.ParentID == "" {
			categories = append(categories, c)
		}
	}
	return categories, nil
}

// Update updates an existing category.
func (r *InMemoryCategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[category.ID]; !exists {
		return entity.ErrCategoryNotFound
	}

	r.categories[category.ID] = category
	return nil
}

// Delete removes a category by its ID.
func (r *InMemoryCategoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[id]; !exists {
		return entity.ErrCategoryNotFound
	}

	delete(r.categories, id)
	return nil
}

// GetActiveCategories retrieves all active categories.
func (r *InMemoryCategoryRepository) GetActiveCategories(ctx context.Context) ([]*entity.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]*entity.Category, 0)
	for _, c := range r.categories {
		if c.Active {
			categories = append(categories, c)
		}
	}
	return categories, nil
}
