package repository

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// CustomerRepository defines the interface for customer persistence operations.
type CustomerRepository interface {
	// Create stores a new customer.
	Create(ctx context.Context, customer *entity.Customer) error
	// GetByID retrieves a customer by its ID.
	GetByID(ctx context.Context, id string) (*entity.Customer, error)
	// GetByEmail retrieves a customer by email.
	GetByEmail(ctx context.Context, email string) (*entity.Customer, error)
	// GetAll retrieves all customers with pagination.
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Customer, error)
	// Update updates an existing customer.
	Update(ctx context.Context, customer *entity.Customer) error
	// Delete removes a customer by its ID.
	Delete(ctx context.Context, id string) error
	// Search searches customers by name or email.
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Customer, error)
	// Count returns the total number of customers.
	Count(ctx context.Context) (int64, error)
}
