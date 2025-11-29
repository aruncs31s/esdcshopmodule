# ESDC Shop Module

A comprehensive shop module built in Go following Domain-Driven Design (DDD) architecture and SOLID principles.

## Features

- **Product Management**: Create, read, update, delete products with stock management
- **Category Management**: Hierarchical category system with parent-child relationships
- **Customer Management**: Customer registration, profile management, and address handling
- **Shopping Cart**: Add/remove items, update quantities, persistent cart per customer
- **Order Management**: Create orders from cart, order status workflow, order history
- **RESTful API**: Clean HTTP API with JSON responses

## Architecture

This project follows Domain-Driven Design (DDD) architecture with clear separation of concerns:

```
├── cmd/
│   └── api/                    # Application entry point
├── internal/
│   ├── domain/                 # Domain Layer (Core Business Logic)
│   │   ├── entity/             # Domain Entities (Product, Order, Cart, Customer, Category)
│   │   ├── valueobject/        # Value Objects (Money, Address)
│   │   ├── repository/         # Repository Interfaces
│   │   └── service/            # Domain Services
│   ├── application/            # Application Layer
│   │   ├── dto/                # Data Transfer Objects
│   │   └── usecase/            # Use Cases (Application Services)
│   ├── infrastructure/         # Infrastructure Layer
│   │   ├── config/             # Configuration
│   │   └── persistence/        # Repository Implementations
│   └── interface/              # Interface Layer
│       └── http/
│           ├── handler/        # HTTP Handlers
│           ├── middleware/     # HTTP Middleware
│           └── router/         # Route Configuration
└── pkg/                        # Shared Packages
    ├── errors/                 # Error Types
    └── validator/              # Validation Utilities
```

## SOLID Principles Applied

- **Single Responsibility Principle (SRP)**: Each class/struct has a single, well-defined responsibility
- **Open/Closed Principle (OCP)**: Entities are open for extension but closed for modification
- **Liskov Substitution Principle (LSP)**: Repository interfaces allow different implementations
- **Interface Segregation Principle (ISP)**: Small, focused interfaces for repositories
- **Dependency Inversion Principle (DIP)**: High-level modules depend on abstractions (interfaces)

## Getting Started

### Prerequisites

- Go 1.22 or later

### Installation

```bash
# Clone the repository
git clone https://github.com/aruncs31s/esdcshopmodule.git
cd esdcshopmodule

# Install dependencies
go mod download

# Build the application
go build -o shop-api ./cmd/api

# Run the application
./shop-api
```

### Running with Go

```bash
go run ./cmd/api
```

The server will start on `http://localhost:8080`

## Configuration

Configuration is done via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_HOST | 0.0.0.0 | Server host address |
| SERVER_PORT | 8080 | Server port |
| APP_NAME | ESDC Shop Module | Application name |
| APP_ENV | development | Environment (development/production) |
| APP_DEBUG | true | Debug mode |
| APP_VERSION | 1.0.0 | Application version |

## API Endpoints

### Health

- `GET /health` - Health check
- `GET /ready` - Readiness check

### Products

- `POST /api/v1/products` - Create product
- `GET /api/v1/products` - List products (with pagination)
- `GET /api/v1/products/{id}` - Get product by ID
- `PUT /api/v1/products/{id}` - Update product
- `DELETE /api/v1/products/{id}` - Delete product
- `PATCH /api/v1/products/{id}/stock` - Update stock
- `GET /api/v1/products/search?q={query}` - Search products
- `GET /api/v1/products/available` - Get available products

### Categories

- `POST /api/v1/categories` - Create category
- `GET /api/v1/categories` - List all categories
- `GET /api/v1/categories/root` - List root categories
- `GET /api/v1/categories/{id}` - Get category by ID
- `PUT /api/v1/categories/{id}` - Update category
- `DELETE /api/v1/categories/{id}` - Delete category
- `GET /api/v1/categories/{id}/subcategories` - Get subcategories
- `GET /api/v1/categories/{id}/products` - Get products by category

### Customers

- `POST /api/v1/customers` - Create customer
- `GET /api/v1/customers` - List customers (with pagination)
- `GET /api/v1/customers/{id}` - Get customer by ID
- `PUT /api/v1/customers/{id}` - Update customer
- `DELETE /api/v1/customers/{id}` - Delete customer
- `PUT /api/v1/customers/{id}/shipping-address` - Update shipping address
- `PUT /api/v1/customers/{id}/billing-address` - Update billing address
- `GET /api/v1/customers/search?q={query}` - Search customers

### Cart

- `GET /api/v1/customers/{id}/cart` - Get cart
- `POST /api/v1/customers/{id}/cart/items` - Add item to cart
- `PUT /api/v1/customers/{id}/cart/items/{productId}` - Update item quantity
- `DELETE /api/v1/customers/{id}/cart/items/{productId}` - Remove item from cart
- `DELETE /api/v1/customers/{id}/cart` - Clear cart

### Orders

- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders` - List orders (with pagination)
- `GET /api/v1/orders/{id}` - Get order by ID
- `PATCH /api/v1/orders/{id}/status` - Update order status
- `POST /api/v1/orders/{id}/cancel` - Cancel order
- `GET /api/v1/customers/{id}/orders` - Get customer orders
- `POST /api/v1/customers/{id}/orders` - Create order from cart

## Example API Usage

### Create a Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wireless Mouse",
    "description": "Ergonomic wireless mouse",
    "price": 2999,
    "currency": "USD",
    "stock": 100,
    "active": true
  }'
```

### Create a Customer

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'
```

### Add Item to Cart

```bash
curl -X POST http://localhost:8080/api/v1/customers/{customer_id}/cart/items \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "{product_id}",
    "quantity": 2
  }'
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests with coverage
go test ./... -cover
```

## Project Structure Benefits

1. **Testability**: Each layer can be tested in isolation
2. **Maintainability**: Clear boundaries between concerns
3. **Flexibility**: Easy to swap implementations (e.g., database)
4. **Scalability**: Can grow the application without major refactoring

## License

MIT License - see [LICENSE](LICENSE) for details
