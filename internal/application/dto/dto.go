package dto

import "time"

// ProductRequest represents a request to create or update a product.
type ProductRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Description string   `json:"description" validate:"max=2000"`
	Price       int64    `json:"price" validate:"required,min=0"`
	Currency    string   `json:"currency" validate:"required,oneof=USD EUR GBP INR"`
	Stock       int      `json:"stock" validate:"min=0"`
	CategoryID  string   `json:"category_id"`
	Images      []string `json:"images"`
	Active      bool     `json:"active"`
}

// ProductResponse represents a product in API responses.
type ProductResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	Currency    string   `json:"currency"`
	Stock       int      `json:"stock"`
	CategoryID  string   `json:"category_id,omitempty"`
	Images      []string `json:"images"`
	Active      bool     `json:"active"`
	Available   bool     `json:"available"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// ProductListResponse represents a paginated list of products.
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Limit      int               `json:"limit"`
	Offset     int               `json:"offset"`
	HasMore    bool              `json:"has_more"`
}

// UpdateStockRequest represents a request to update product stock.
type UpdateStockRequest struct {
	Stock int `json:"stock" validate:"required,min=0"`
}

// SearchProductRequest represents a product search request.
type SearchProductRequest struct {
	Query  string `json:"query" validate:"required,min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Offset int    `json:"offset" validate:"min=0"`
}

// CategoryRequest represents a request to create or update a category.
type CategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	ParentID    string `json:"parent_id"`
}

// CategoryResponse represents a category in API responses.
type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id,omitempty"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CategoryListResponse represents a list of categories.
type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int                `json:"total"`
}

// CartItemRequest represents a request to add or update a cart item.
type CartItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// CartItemResponse represents a cart item in API responses.
type CartItemResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name,omitempty"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`
	Currency    string `json:"currency"`
	TotalPrice  int64  `json:"total_price"`
	AddedAt     string `json:"added_at"`
}

// CartResponse represents a cart in API responses.
type CartResponse struct {
	ID         string             `json:"id"`
	CustomerID string             `json:"customer_id"`
	Items      []CartItemResponse `json:"items"`
	ItemCount  int                `json:"item_count"`
	Total      int64              `json:"total"`
	Currency   string             `json:"currency"`
	CreatedAt  string             `json:"created_at"`
	UpdatedAt  string             `json:"updated_at"`
}

// OrderItemRequest represents an item in an order request.
type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// CreateOrderRequest represents a request to create an order.
type CreateOrderRequest struct {
	CustomerID      string             `json:"customer_id" validate:"required"`
	Items           []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
	ShippingAddress AddressRequest     `json:"shipping_address" validate:"required"`
	BillingAddress  AddressRequest     `json:"billing_address" validate:"required"`
	Notes           string             `json:"notes" validate:"max=1000"`
}

// CreateOrderFromCartRequest represents a request to create an order from cart.
type CreateOrderFromCartRequest struct {
	ShippingAddress AddressRequest `json:"shipping_address" validate:"required"`
	BillingAddress  AddressRequest `json:"billing_address" validate:"required"`
	Notes           string         `json:"notes" validate:"max=1000"`
}

// AddressRequest represents an address in requests.
type AddressRequest struct {
	Street     string `json:"street" validate:"required,min=1,max=255"`
	City       string `json:"city" validate:"required,min=1,max=100"`
	State      string `json:"state" validate:"max=100"`
	PostalCode string `json:"postal_code" validate:"required,min=1,max=20"`
	Country    string `json:"country" validate:"required,min=1,max=100"`
}

// AddressResponse represents an address in API responses.
type AddressResponse struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// OrderItemResponse represents an order item in API responses.
type OrderItemResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`
	Currency    string `json:"currency"`
	TotalPrice  int64  `json:"total_price"`
}

// OrderResponse represents an order in API responses.
type OrderResponse struct {
	ID              string              `json:"id"`
	CustomerID      string              `json:"customer_id"`
	Items           []OrderItemResponse `json:"items"`
	ItemCount       int                 `json:"item_count"`
	TotalAmount     int64               `json:"total_amount"`
	Currency        string              `json:"currency"`
	Status          string              `json:"status"`
	ShippingAddress AddressResponse     `json:"shipping_address"`
	BillingAddress  AddressResponse     `json:"billing_address"`
	Notes           string              `json:"notes,omitempty"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
}

// OrderListResponse represents a paginated list of orders.
type OrderListResponse struct {
	Orders  []OrderResponse `json:"orders"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"has_more"`
}

// UpdateOrderStatusRequest represents a request to update order status.
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=confirmed processing shipped delivered cancelled refunded"`
}

// CustomerRequest represents a request to create or update a customer.
type CustomerRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Phone     string `json:"phone" validate:"max=20"`
}

// CustomerResponse represents a customer in API responses.
type CustomerResponse struct {
	ID              string           `json:"id"`
	Email           string           `json:"email"`
	FirstName       string           `json:"first_name"`
	LastName        string           `json:"last_name"`
	Phone           string           `json:"phone,omitempty"`
	FullName        string           `json:"full_name"`
	ShippingAddress *AddressResponse `json:"shipping_address,omitempty"`
	BillingAddress  *AddressResponse `json:"billing_address,omitempty"`
	Active          bool             `json:"active"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
}

// CustomerListResponse represents a paginated list of customers.
type CustomerListResponse struct {
	Customers []CustomerResponse `json:"customers"`
	Total     int64              `json:"total"`
	Limit     int                `json:"limit"`
	Offset    int                `json:"offset"`
	HasMore   bool               `json:"has_more"`
}

// PaginationRequest represents pagination parameters.
type PaginationRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response.
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}
