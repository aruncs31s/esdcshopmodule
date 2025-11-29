package router

import (
	"net/http"

	"github.com/aruncs31s/esdcshopmodule/internal/interface/http/handler"
	"github.com/aruncs31s/esdcshopmodule/internal/interface/http/middleware"
)

// Router holds all HTTP handlers and configures routes.
type Router struct {
	ProductHandler  *handler.ProductHandler
	OrderHandler    *handler.OrderHandler
	CartHandler     *handler.CartHandler
	CategoryHandler *handler.CategoryHandler
	CustomerHandler *handler.CustomerHandler
	HealthHandler   *handler.HealthHandler
}

// NewRouter creates a new Router.
func NewRouter(
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
	cartHandler *handler.CartHandler,
	categoryHandler *handler.CategoryHandler,
	customerHandler *handler.CustomerHandler,
	healthHandler *handler.HealthHandler,
) *Router {
	return &Router{
		ProductHandler:  productHandler,
		OrderHandler:    orderHandler,
		CartHandler:     cartHandler,
		CategoryHandler: categoryHandler,
		CustomerHandler: customerHandler,
		HealthHandler:   healthHandler,
	}
}

// Setup configures all routes and returns the HTTP handler.
func (r *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// Health routes
	mux.HandleFunc("GET /health", r.HealthHandler.Health)
	mux.HandleFunc("GET /ready", r.HealthHandler.Ready)

	// Product routes
	mux.HandleFunc("POST /api/v1/products", r.ProductHandler.CreateProduct)
	mux.HandleFunc("GET /api/v1/products", r.ProductHandler.GetProducts)
	mux.HandleFunc("GET /api/v1/products/search", r.ProductHandler.SearchProducts)
	mux.HandleFunc("GET /api/v1/products/available", r.ProductHandler.GetAvailableProducts)
	mux.HandleFunc("GET /api/v1/products/{id}", r.ProductHandler.GetProduct)
	mux.HandleFunc("PUT /api/v1/products/{id}", r.ProductHandler.UpdateProduct)
	mux.HandleFunc("DELETE /api/v1/products/{id}", r.ProductHandler.DeleteProduct)
	mux.HandleFunc("PATCH /api/v1/products/{id}/stock", r.ProductHandler.UpdateStock)

	// Category routes
	mux.HandleFunc("POST /api/v1/categories", r.CategoryHandler.CreateCategory)
	mux.HandleFunc("GET /api/v1/categories", r.CategoryHandler.GetCategories)
	mux.HandleFunc("GET /api/v1/categories/root", r.CategoryHandler.GetRootCategories)
	mux.HandleFunc("GET /api/v1/categories/{id}", r.CategoryHandler.GetCategory)
	mux.HandleFunc("PUT /api/v1/categories/{id}", r.CategoryHandler.UpdateCategory)
	mux.HandleFunc("DELETE /api/v1/categories/{id}", r.CategoryHandler.DeleteCategory)
	mux.HandleFunc("GET /api/v1/categories/{id}/subcategories", r.CategoryHandler.GetSubcategories)
	mux.HandleFunc("GET /api/v1/categories/{id}/products", r.ProductHandler.GetProductsByCategory)

	// Customer routes
	mux.HandleFunc("POST /api/v1/customers", r.CustomerHandler.CreateCustomer)
	mux.HandleFunc("GET /api/v1/customers", r.CustomerHandler.GetCustomers)
	mux.HandleFunc("GET /api/v1/customers/search", r.CustomerHandler.SearchCustomers)
	mux.HandleFunc("GET /api/v1/customers/{id}", r.CustomerHandler.GetCustomer)
	mux.HandleFunc("PUT /api/v1/customers/{id}", r.CustomerHandler.UpdateCustomer)
	mux.HandleFunc("DELETE /api/v1/customers/{id}", r.CustomerHandler.DeleteCustomer)
	mux.HandleFunc("PUT /api/v1/customers/{id}/shipping-address", r.CustomerHandler.UpdateShippingAddress)
	mux.HandleFunc("PUT /api/v1/customers/{id}/billing-address", r.CustomerHandler.UpdateBillingAddress)

	// Cart routes
	mux.HandleFunc("GET /api/v1/customers/{id}/cart", r.CartHandler.GetCart)
	mux.HandleFunc("POST /api/v1/customers/{id}/cart/items", r.CartHandler.AddItem)
	mux.HandleFunc("PUT /api/v1/customers/{id}/cart/items/{productId}", r.CartHandler.UpdateItemQuantity)
	mux.HandleFunc("DELETE /api/v1/customers/{id}/cart/items/{productId}", r.CartHandler.RemoveItem)
	mux.HandleFunc("DELETE /api/v1/customers/{id}/cart", r.CartHandler.ClearCart)

	// Order routes
	mux.HandleFunc("POST /api/v1/orders", r.OrderHandler.CreateOrder)
	mux.HandleFunc("GET /api/v1/orders", r.OrderHandler.GetOrders)
	mux.HandleFunc("GET /api/v1/orders/{id}", r.OrderHandler.GetOrder)
	mux.HandleFunc("PATCH /api/v1/orders/{id}/status", r.OrderHandler.UpdateOrderStatus)
	mux.HandleFunc("POST /api/v1/orders/{id}/cancel", r.OrderHandler.CancelOrder)
	mux.HandleFunc("GET /api/v1/customers/{id}/orders", r.OrderHandler.GetCustomerOrders)
	mux.HandleFunc("POST /api/v1/customers/{id}/orders", r.OrderHandler.CreateOrderFromCart)

	// Apply middleware
	return middleware.Chain(
		mux,
		middleware.Recovery,
		middleware.Logging,
		middleware.CORS,
		middleware.ContentType,
	)
}
