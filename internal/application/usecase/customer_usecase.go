package usecase

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// CustomerUseCase handles customer-related business logic.
type CustomerUseCase struct {
	customerRepo repository.CustomerRepository
	idGenerator  IDGenerator
}

// NewCustomerUseCase creates a new CustomerUseCase.
func NewCustomerUseCase(
	customerRepo repository.CustomerRepository,
	idGenerator IDGenerator,
) *CustomerUseCase {
	return &CustomerUseCase{
		customerRepo: customerRepo,
		idGenerator:  idGenerator,
	}
}

// CreateCustomer creates a new customer.
func (uc *CustomerUseCase) CreateCustomer(ctx context.Context, req dto.CustomerRequest) (*dto.CustomerResponse, error) {
	// Check if email already exists
	existing, _ := uc.customerRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, entity.ErrInvalidCustomerEmail
	}

	// Generate ID
	id := uc.idGenerator.Generate()

	// Create customer entity
	customer, err := entity.NewCustomer(id, req.Email, req.FirstName, req.LastName, req.Phone)
	if err != nil {
		return nil, err
	}

	// Save customer
	if err := uc.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return uc.toCustomerResponse(customer), nil
}

// GetCustomer retrieves a customer by ID.
func (uc *CustomerUseCase) GetCustomer(ctx context.Context, id string) (*dto.CustomerResponse, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toCustomerResponse(customer), nil
}

// GetCustomerByEmail retrieves a customer by email.
func (uc *CustomerUseCase) GetCustomerByEmail(ctx context.Context, email string) (*dto.CustomerResponse, error) {
	customer, err := uc.customerRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return uc.toCustomerResponse(customer), nil
}

// GetCustomers retrieves all customers with pagination.
func (uc *CustomerUseCase) GetCustomers(ctx context.Context, limit, offset int) (*dto.CustomerListResponse, error) {
	customers, err := uc.customerRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.customerRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	customerResponses := make([]dto.CustomerResponse, len(customers))
	for i, c := range customers {
		customerResponses[i] = *uc.toCustomerResponse(c)
	}

	return &dto.CustomerListResponse{
		Customers: customerResponses,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
		HasMore:   int64(offset+len(customers)) < total,
	}, nil
}

// UpdateCustomer updates an existing customer.
func (uc *CustomerUseCase) UpdateCustomer(ctx context.Context, id string, req dto.CustomerRequest) (*dto.CustomerResponse, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed to an existing one
	if customer.Email != req.Email {
		existing, _ := uc.customerRepo.GetByEmail(ctx, req.Email)
		if existing != nil {
			return nil, entity.ErrInvalidCustomerEmail
		}
		if err := customer.UpdateEmail(req.Email); err != nil {
			return nil, err
		}
	}

	// Update other fields
	customer.FirstName = req.FirstName
	customer.LastName = req.LastName
	customer.Phone = req.Phone

	if err := uc.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return uc.toCustomerResponse(customer), nil
}

// UpdateShippingAddress updates a customer's shipping address.
func (uc *CustomerUseCase) UpdateShippingAddress(ctx context.Context, id string, req dto.AddressRequest) (*dto.CustomerResponse, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	addr, err := valueobject.NewAddress(req.Street, req.City, req.State, req.PostalCode, req.Country)
	if err != nil {
		return nil, err
	}

	customer.SetShippingAddress(addr)

	if err := uc.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return uc.toCustomerResponse(customer), nil
}

// UpdateBillingAddress updates a customer's billing address.
func (uc *CustomerUseCase) UpdateBillingAddress(ctx context.Context, id string, req dto.AddressRequest) (*dto.CustomerResponse, error) {
	customer, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	addr, err := valueobject.NewAddress(req.Street, req.City, req.State, req.PostalCode, req.Country)
	if err != nil {
		return nil, err
	}

	customer.SetBillingAddress(addr)

	if err := uc.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return uc.toCustomerResponse(customer), nil
}

// DeleteCustomer deletes a customer.
func (uc *CustomerUseCase) DeleteCustomer(ctx context.Context, id string) error {
	_, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.customerRepo.Delete(ctx, id)
}

// SearchCustomers searches customers by name or email.
func (uc *CustomerUseCase) SearchCustomers(ctx context.Context, query string, limit, offset int) (*dto.CustomerListResponse, error) {
	customers, err := uc.customerRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	customerResponses := make([]dto.CustomerResponse, len(customers))
	for i, c := range customers {
		customerResponses[i] = *uc.toCustomerResponse(c)
	}

	return &dto.CustomerListResponse{
		Customers: customerResponses,
		Total:     int64(len(customers)),
		Limit:     limit,
		Offset:    offset,
		HasMore:   false,
	}, nil
}

// toCustomerResponse converts a customer entity to a response DTO.
func (uc *CustomerUseCase) toCustomerResponse(customer *entity.Customer) *dto.CustomerResponse {
	resp := &dto.CustomerResponse{
		ID:        customer.ID,
		Email:     customer.Email,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Phone:     customer.Phone,
		FullName:  customer.FullName(),
		Active:    customer.Active,
		CreatedAt: customer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: customer.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if customer.ShippingAddress.Street != "" {
		resp.ShippingAddress = &dto.AddressResponse{
			Street:     customer.ShippingAddress.Street,
			City:       customer.ShippingAddress.City,
			State:      customer.ShippingAddress.State,
			PostalCode: customer.ShippingAddress.PostalCode,
			Country:    customer.ShippingAddress.Country,
		}
	}

	if customer.BillingAddress.Street != "" {
		resp.BillingAddress = &dto.AddressResponse{
			Street:     customer.BillingAddress.Street,
			City:       customer.BillingAddress.City,
			State:      customer.BillingAddress.State,
			PostalCode: customer.BillingAddress.PostalCode,
			Country:    customer.BillingAddress.Country,
		}
	}

	return resp
}
