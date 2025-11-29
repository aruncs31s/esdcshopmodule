package persistence

import (
	"context"
	"sync"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
)

// InMemoryOrderRepository is an in-memory implementation of OrderRepository.
type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*entity.Order
}

// NewInMemoryOrderRepository creates a new InMemoryOrderRepository.
func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*entity.Order),
	}
}

// Ensure InMemoryOrderRepository implements OrderRepository.
var _ repository.OrderRepository = (*InMemoryOrderRepository)(nil)

// Create stores a new order.
func (r *InMemoryOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; exists {
		return entity.ErrInvalidOrderID
	}

	r.orders[order.ID] = order
	return nil
}

// GetByID retrieves an order by its ID.
func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, entity.ErrOrderNotFound
	}
	return order, nil
}

// GetByCustomerID retrieves orders by customer ID.
func (r *InMemoryOrderRepository) GetByCustomerID(ctx context.Context, customerID string, limit, offset int) ([]*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*entity.Order, 0)
	for _, o := range r.orders {
		if o.CustomerID == customerID {
			orders = append(orders, o)
		}
	}

	// Apply pagination
	if offset >= len(orders) {
		return []*entity.Order{}, nil
	}

	end := offset + limit
	if end > len(orders) {
		end = len(orders)
	}

	return orders[offset:end], nil
}

// GetAll retrieves all orders with pagination.
func (r *InMemoryOrderRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*entity.Order, 0, len(r.orders))
	for _, o := range r.orders {
		orders = append(orders, o)
	}

	// Apply pagination
	if offset >= len(orders) {
		return []*entity.Order{}, nil
	}

	end := offset + limit
	if end > len(orders) {
		end = len(orders)
	}

	return orders[offset:end], nil
}

// GetByStatus retrieves orders by status.
func (r *InMemoryOrderRepository) GetByStatus(ctx context.Context, status entity.OrderStatus, limit, offset int) ([]*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*entity.Order, 0)
	for _, o := range r.orders {
		if o.Status == status {
			orders = append(orders, o)
		}
	}

	// Apply pagination
	if offset >= len(orders) {
		return []*entity.Order{}, nil
	}

	end := offset + limit
	if end > len(orders) {
		end = len(orders)
	}

	return orders[offset:end], nil
}

// Update updates an existing order.
func (r *InMemoryOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return entity.ErrOrderNotFound
	}

	r.orders[order.ID] = order
	return nil
}

// UpdateStatus updates the status of an order.
func (r *InMemoryOrderRepository) UpdateStatus(ctx context.Context, id string, status entity.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.orders[id]
	if !exists {
		return entity.ErrOrderNotFound
	}

	return order.UpdateStatus(status)
}

// Delete removes an order by its ID.
func (r *InMemoryOrderRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[id]; !exists {
		return entity.ErrOrderNotFound
	}

	delete(r.orders, id)
	return nil
}

// Count returns the total number of orders.
func (r *InMemoryOrderRepository) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return int64(len(r.orders)), nil
}

// CountByCustomer returns the total number of orders for a customer.
func (r *InMemoryOrderRepository) CountByCustomer(ctx context.Context, customerID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := int64(0)
	for _, o := range r.orders {
		if o.CustomerID == customerID {
			count++
		}
	}
	return count, nil
}
