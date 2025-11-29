package persistence

import (
	"context"
	"sync"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
)

// InMemoryCartRepository is an in-memory implementation of CartRepository.
type InMemoryCartRepository struct {
	mu    sync.RWMutex
	carts map[string]*entity.Cart
}

// NewInMemoryCartRepository creates a new InMemoryCartRepository.
func NewInMemoryCartRepository() *InMemoryCartRepository {
	return &InMemoryCartRepository{
		carts: make(map[string]*entity.Cart),
	}
}

// Ensure InMemoryCartRepository implements CartRepository.
var _ repository.CartRepository = (*InMemoryCartRepository)(nil)

// Create stores a new cart.
func (r *InMemoryCartRepository) Create(ctx context.Context, cart *entity.Cart) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.carts[cart.ID]; exists {
		return entity.ErrInvalidCartID
	}

	r.carts[cart.ID] = cart
	return nil
}

// GetByID retrieves a cart by its ID.
func (r *InMemoryCartRepository) GetByID(ctx context.Context, id string) (*entity.Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cart, exists := r.carts[id]
	if !exists {
		return nil, entity.ErrCartNotFound
	}
	return cart, nil
}

// GetByCustomerID retrieves a cart by customer ID.
func (r *InMemoryCartRepository) GetByCustomerID(ctx context.Context, customerID string) (*entity.Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, cart := range r.carts {
		if cart.CustomerID == customerID {
			return cart, nil
		}
	}
	return nil, entity.ErrCartNotFound
}

// Update updates an existing cart.
func (r *InMemoryCartRepository) Update(ctx context.Context, cart *entity.Cart) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.carts[cart.ID]; !exists {
		return entity.ErrCartNotFound
	}

	r.carts[cart.ID] = cart
	return nil
}

// Delete removes a cart by its ID.
func (r *InMemoryCartRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.carts[id]; !exists {
		return entity.ErrCartNotFound
	}

	delete(r.carts, id)
	return nil
}

// DeleteByCustomerID removes a cart by customer ID.
func (r *InMemoryCartRepository) DeleteByCustomerID(ctx context.Context, customerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, cart := range r.carts {
		if cart.CustomerID == customerID {
			delete(r.carts, id)
			return nil
		}
	}
	return entity.ErrCartNotFound
}
