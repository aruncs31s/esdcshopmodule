package entity

import (
	"time"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// CartItem represents an item in a shopping cart.
type CartItem struct {
	ProductID  string
	Quantity   int
	UnitPrice  valueobject.Money
	TotalPrice valueobject.Money
	AddedAt    time.Time
}

// Cart represents a shopping cart entity.
type Cart struct {
	ID         string
	CustomerID string
	Items      []CartItem
	Total      valueobject.Money
	Currency   valueobject.Currency
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewCart creates a new Cart with validation.
func NewCart(id, customerID string, currency valueobject.Currency) (*Cart, error) {
	if id == "" {
		return nil, ErrInvalidCartID
	}
	if customerID == "" {
		return nil, ErrInvalidCustomerID
	}

	now := time.Now()
	return &Cart{
		ID:         id,
		CustomerID: customerID,
		Items:      []CartItem{},
		Total:      valueobject.Money{Amount: 0, Currency: currency},
		Currency:   currency,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// AddItem adds an item to the cart.
func (c *Cart) AddItem(productID string, quantity int, unitPrice valueobject.Money) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	// Check if item already exists
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity += quantity
			c.Items[i].TotalPrice = unitPrice.Multiply(c.Items[i].Quantity)
			c.recalculateTotal()
			c.UpdatedAt = time.Now()
			return nil
		}
	}

	// Add new item
	newItem := CartItem{
		ProductID:  productID,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: unitPrice.Multiply(quantity),
		AddedAt:    time.Now(),
	}
	c.Items = append(c.Items, newItem)
	c.recalculateTotal()
	c.UpdatedAt = time.Now()
	return nil
}

// UpdateItemQuantity updates the quantity of an item in the cart.
func (c *Cart) UpdateItemQuantity(productID string, quantity int) error {
	if quantity <= 0 {
		return c.RemoveItem(productID)
	}

	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity = quantity
			c.Items[i].TotalPrice = item.UnitPrice.Multiply(quantity)
			c.recalculateTotal()
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrProductNotFound
}

// RemoveItem removes an item from the cart.
func (c *Cart) RemoveItem(productID string) error {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.recalculateTotal()
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrProductNotFound
}

// Clear removes all items from the cart.
func (c *Cart) Clear() {
	c.Items = []CartItem{}
	c.Total = valueobject.Money{Amount: 0, Currency: c.Currency}
	c.UpdatedAt = time.Now()
}

// GetItem returns an item from the cart by product ID.
func (c *Cart) GetItem(productID string) (*CartItem, error) {
	for _, item := range c.Items {
		if item.ProductID == productID {
			return &item, nil
		}
	}
	return nil, ErrProductNotFound
}

// IsEmpty checks if the cart is empty.
func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

// GetItemCount returns the total number of items in the cart.
func (c *Cart) GetItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

// recalculateTotal recalculates the cart total.
func (c *Cart) recalculateTotal() {
	total := valueobject.Money{Amount: 0, Currency: c.Currency}
	for _, item := range c.Items {
		total, _ = total.Add(item.TotalPrice)
	}
	c.Total = total
}
