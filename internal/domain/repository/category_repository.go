package repository

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// CategoryRepository defines the interface for category persistence operations.
type CategoryRepository interface {
	// Create stores a new category.
	Create(ctx context.Context, category *entity.Category) error
	// GetByID retrieves a category by its ID.
	GetByID(ctx context.Context, id string) (*entity.Category, error)
	// GetAll retrieves all categories.
	GetAll(ctx context.Context) ([]*entity.Category, error)
	// GetByParentID retrieves categories by parent ID.
	GetByParentID(ctx context.Context, parentID string) ([]*entity.Category, error)
	// GetRootCategories retrieves all root categories (no parent).
	GetRootCategories(ctx context.Context) ([]*entity.Category, error)
	// Update updates an existing category.
	Update(ctx context.Context, category *entity.Category) error
	// Delete removes a category by its ID.
	Delete(ctx context.Context, id string) error
	// GetActiveCategories retrieves all active categories.
	GetActiveCategories(ctx context.Context) ([]*entity.Category, error)
}
