package entity

import "time"

// Category represents a product category entity.
type Category struct {
	ID          string
	Name        string
	Description string
	ParentID    string // For hierarchical categories
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewCategory creates a new Category with validation.
func NewCategory(id, name, description, parentID string) (*Category, error) {
	if id == "" {
		return nil, ErrInvalidCategoryID
	}
	if name == "" {
		return nil, ErrInvalidCategoryName
	}

	now := time.Now()
	return &Category{
		ID:          id,
		Name:        name,
		Description: description,
		ParentID:    parentID,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update updates the category details.
func (c *Category) Update(name, description string) error {
	if name == "" {
		return ErrInvalidCategoryName
	}
	c.Name = name
	c.Description = description
	c.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the category.
func (c *Category) Deactivate() {
	c.Active = false
	c.UpdatedAt = time.Now()
}

// Activate activates the category.
func (c *Category) Activate() {
	c.Active = true
	c.UpdatedAt = time.Now()
}

// HasParent checks if the category has a parent category.
func (c *Category) HasParent() bool {
	return c.ParentID != ""
}
