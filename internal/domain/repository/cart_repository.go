package repository

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// CartRepository defines the interface for cart persistence operations.
type CartRepository interface {
	// Create stores a new cart.
	Create(ctx context.Context, cart *entity.Cart) error
	// GetByID retrieves a cart by its ID.
	GetByID(ctx context.Context, id string) (*entity.Cart, error)
	// GetByCustomerID retrieves a cart by customer ID.
	GetByCustomerID(ctx context.Context, customerID string) (*entity.Cart, error)
	// Update updates an existing cart.
	Update(ctx context.Context, cart *entity.Cart) error
	// Delete removes a cart by its ID.
	Delete(ctx context.Context, id string) error
	// DeleteByCustomerID removes a cart by customer ID.
	DeleteByCustomerID(ctx context.Context, customerID string) error
}
