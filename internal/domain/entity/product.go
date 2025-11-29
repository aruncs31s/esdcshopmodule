package entity

import (
	"time"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// Product represents a product entity in the shop domain.
type Product struct {
	ID          string
	Name        string
	Description string
	Price       valueobject.Money
	Stock       int
	CategoryID  string
	Images      []string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProduct creates a new Product with validation.
func NewProduct(id, name, description string, price valueobject.Money, stock int, categoryID string) (*Product, error) {
	if id == "" {
		return nil, ErrInvalidProductID
	}
	if name == "" {
		return nil, ErrInvalidProductName
	}
	if stock < 0 {
		return nil, ErrInvalidStock
	}

	now := time.Now()
	return &Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CategoryID:  categoryID,
		Images:      []string{},
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateStock updates the product stock.
func (p *Product) UpdateStock(quantity int) error {
	if quantity < 0 {
		return ErrInvalidStock
	}
	p.Stock = quantity
	p.UpdatedAt = time.Now()
	return nil
}

// DecreaseStock decreases the product stock by the given quantity.
func (p *Product) DecreaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if p.Stock < quantity {
		return ErrInsufficientStock
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}

// IncreaseStock increases the product stock by the given quantity.
func (p *Product) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	p.Stock += quantity
	p.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the product.
func (p *Product) Deactivate() {
	p.Active = false
	p.UpdatedAt = time.Now()
}

// Activate activates the product.
func (p *Product) Activate() {
	p.Active = true
	p.UpdatedAt = time.Now()
}

// AddImage adds an image URL to the product.
func (p *Product) AddImage(imageURL string) {
	p.Images = append(p.Images, imageURL)
	p.UpdatedAt = time.Now()
}

// IsAvailable checks if the product is available for purchase.
func (p *Product) IsAvailable() bool {
	return p.Active && p.Stock > 0
}
