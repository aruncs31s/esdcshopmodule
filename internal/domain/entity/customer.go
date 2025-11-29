package entity

import (
	"regexp"
	"time"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// Customer represents a customer entity in the shop domain.
type Customer struct {
	ID              string
	Email           string
	FirstName       string
	LastName        string
	Phone           string
	ShippingAddress valueobject.Address
	BillingAddress  valueobject.Address
	Active          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewCustomer creates a new Customer with validation.
func NewCustomer(id, email, firstName, lastName, phone string) (*Customer, error) {
	if id == "" {
		return nil, ErrInvalidCustomerID
	}
	if !isValidEmail(email) {
		return nil, ErrInvalidCustomerEmail
	}

	now := time.Now()
	return &Customer{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// FullName returns the customer's full name.
func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}

// UpdateEmail updates the customer's email.
func (c *Customer) UpdateEmail(email string) error {
	if !isValidEmail(email) {
		return ErrInvalidCustomerEmail
	}
	c.Email = email
	c.UpdatedAt = time.Now()
	return nil
}

// SetShippingAddress sets the customer's shipping address.
func (c *Customer) SetShippingAddress(addr valueobject.Address) {
	c.ShippingAddress = addr
	c.UpdatedAt = time.Now()
}

// SetBillingAddress sets the customer's billing address.
func (c *Customer) SetBillingAddress(addr valueobject.Address) {
	c.BillingAddress = addr
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the customer account.
func (c *Customer) Deactivate() {
	c.Active = false
	c.UpdatedAt = time.Now()
}

// Activate activates the customer account.
func (c *Customer) Activate() {
	c.Active = true
	c.UpdatedAt = time.Now()
}

// isValidEmail validates an email address.
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
