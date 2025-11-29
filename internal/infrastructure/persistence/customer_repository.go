package persistence

import (
	"context"
	"strings"
	"sync"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
)

// InMemoryCustomerRepository is an in-memory implementation of CustomerRepository.
type InMemoryCustomerRepository struct {
	mu        sync.RWMutex
	customers map[string]*entity.Customer
}

// NewInMemoryCustomerRepository creates a new InMemoryCustomerRepository.
func NewInMemoryCustomerRepository() *InMemoryCustomerRepository {
	return &InMemoryCustomerRepository{
		customers: make(map[string]*entity.Customer),
	}
}

// Ensure InMemoryCustomerRepository implements CustomerRepository.
var _ repository.CustomerRepository = (*InMemoryCustomerRepository)(nil)

// Create stores a new customer.
func (r *InMemoryCustomerRepository) Create(ctx context.Context, customer *entity.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.customers[customer.ID]; exists {
		return entity.ErrInvalidCustomerID
	}

	r.customers[customer.ID] = customer
	return nil
}

// GetByID retrieves a customer by its ID.
func (r *InMemoryCustomerRepository) GetByID(ctx context.Context, id string) (*entity.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customer, exists := r.customers[id]
	if !exists {
		return nil, entity.ErrCustomerNotFound
	}
	return customer, nil
}

// GetByEmail retrieves a customer by email.
func (r *InMemoryCustomerRepository) GetByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, customer := range r.customers {
		if customer.Email == email {
			return customer, nil
		}
	}
	return nil, entity.ErrCustomerNotFound
}

// GetAll retrieves all customers with pagination.
func (r *InMemoryCustomerRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customers := make([]*entity.Customer, 0, len(r.customers))
	for _, c := range r.customers {
		customers = append(customers, c)
	}

	// Apply pagination
	if offset >= len(customers) {
		return []*entity.Customer{}, nil
	}

	end := offset + limit
	if end > len(customers) {
		end = len(customers)
	}

	return customers[offset:end], nil
}

// Update updates an existing customer.
func (r *InMemoryCustomerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.customers[customer.ID]; !exists {
		return entity.ErrCustomerNotFound
	}

	r.customers[customer.ID] = customer
	return nil
}

// Delete removes a customer by its ID.
func (r *InMemoryCustomerRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.customers[id]; !exists {
		return entity.ErrCustomerNotFound
	}

	delete(r.customers, id)
	return nil
}

// Search searches customers by name or email.
func (r *InMemoryCustomerRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(query)
	customers := make([]*entity.Customer, 0)

	for _, c := range r.customers {
		if strings.Contains(strings.ToLower(c.Email), query) ||
			strings.Contains(strings.ToLower(c.FirstName), query) ||
			strings.Contains(strings.ToLower(c.LastName), query) {
			customers = append(customers, c)
		}
	}

	// Apply pagination
	if offset >= len(customers) {
		return []*entity.Customer{}, nil
	}

	end := offset + limit
	if end > len(customers) {
		end = len(customers)
	}

	return customers[offset:end], nil
}

// Count returns the total number of customers.
func (r *InMemoryCustomerRepository) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return int64(len(r.customers)), nil
}
