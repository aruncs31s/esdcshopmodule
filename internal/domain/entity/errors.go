package entity

import "errors"

// Product errors
var (
	ErrInvalidProductID   = errors.New("invalid product ID")
	ErrInvalidProductName = errors.New("invalid product name")
	ErrInvalidStock       = errors.New("invalid stock quantity")
	ErrInvalidQuantity    = errors.New("invalid quantity")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrProductNotFound    = errors.New("product not found")
)

// Order errors
var (
	ErrInvalidOrderID     = errors.New("invalid order ID")
	ErrInvalidCustomerID  = errors.New("invalid customer ID")
	ErrEmptyOrder         = errors.New("order cannot be empty")
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status transition")
)

// Category errors
var (
	ErrInvalidCategoryID   = errors.New("invalid category ID")
	ErrInvalidCategoryName = errors.New("invalid category name")
	ErrCategoryNotFound    = errors.New("category not found")
)

// Cart errors
var (
	ErrInvalidCartID  = errors.New("invalid cart ID")
	ErrCartNotFound   = errors.New("cart not found")
	ErrCartItemExists = errors.New("item already exists in cart")
)

// Customer errors
var (
	ErrInvalidCustomerEmail = errors.New("invalid customer email")
	ErrCustomerNotFound     = errors.New("customer not found")
)
