package repository

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
)

// OrderRepository defines the interface for order persistence operations.
type OrderRepository interface {
	// Create stores a new order.
	Create(ctx context.Context, order *entity.Order) error
	// GetByID retrieves an order by its ID.
	GetByID(ctx context.Context, id string) (*entity.Order, error)
	// GetByCustomerID retrieves orders by customer ID.
	GetByCustomerID(ctx context.Context, customerID string, limit, offset int) ([]*entity.Order, error)
	// GetAll retrieves all orders with pagination.
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Order, error)
	// GetByStatus retrieves orders by status.
	GetByStatus(ctx context.Context, status entity.OrderStatus, limit, offset int) ([]*entity.Order, error)
	// Update updates an existing order.
	Update(ctx context.Context, order *entity.Order) error
	// UpdateStatus updates the status of an order.
	UpdateStatus(ctx context.Context, id string, status entity.OrderStatus) error
	// Delete removes an order by its ID.
	Delete(ctx context.Context, id string) error
	// Count returns the total number of orders.
	Count(ctx context.Context) (int64, error)
	// CountByCustomer returns the total number of orders for a customer.
	CountByCustomer(ctx context.Context, customerID string) (int64, error)
}
