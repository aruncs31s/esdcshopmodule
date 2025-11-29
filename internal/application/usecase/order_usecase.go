package usecase

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/application/dto"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/service"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// OrderUseCase handles order-related business logic.
type OrderUseCase struct {
	orderRepo    repository.OrderRepository
	productRepo  repository.ProductRepository
	cartRepo     repository.CartRepository
	customerRepo repository.CustomerRepository
	orderService *service.OrderService
	idGenerator  IDGenerator
}

// NewOrderUseCase creates a new OrderUseCase.
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	cartRepo repository.CartRepository,
	customerRepo repository.CustomerRepository,
	orderService *service.OrderService,
	idGenerator IDGenerator,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		cartRepo:     cartRepo,
		customerRepo: customerRepo,
		orderService: orderService,
		idGenerator:  idGenerator,
	}
}

// CreateOrder creates a new order.
func (uc *OrderUseCase) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// Validate customer exists
	_, err := uc.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, entity.ErrCustomerNotFound
	}

	// Convert address DTOs to value objects
	shippingAddr, err := valueobject.NewAddress(
		req.ShippingAddress.Street,
		req.ShippingAddress.City,
		req.ShippingAddress.State,
		req.ShippingAddress.PostalCode,
		req.ShippingAddress.Country,
	)
	if err != nil {
		return nil, err
	}

	billingAddr, err := valueobject.NewAddress(
		req.BillingAddress.Street,
		req.BillingAddress.City,
		req.BillingAddress.State,
		req.BillingAddress.PostalCode,
		req.BillingAddress.Country,
	)
	if err != nil {
		return nil, err
	}

	// Build order items
	orderItems := make([]entity.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		product, err := uc.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, entity.ErrProductNotFound
		}

		if product.Stock < item.Quantity {
			return nil, entity.ErrInsufficientStock
		}

		totalPrice := product.Price.Multiply(item.Quantity)
		orderItem := entity.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			TotalPrice:  totalPrice,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Generate order ID
	orderID := uc.idGenerator.Generate()

	// Create order
	order, err := entity.NewOrder(orderID, req.CustomerID, orderItems, shippingAddr, billingAddr)
	if err != nil {
		return nil, err
	}
	order.Notes = req.Notes

	// Decrease stock for each product
	for _, item := range req.Items {
		product, _ := uc.productRepo.GetByID(ctx, item.ProductID)
		if err := product.DecreaseStock(item.Quantity); err != nil {
			return nil, err
		}
		if err := uc.productRepo.Update(ctx, product); err != nil {
			return nil, err
		}
	}

	// Save order
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return uc.toOrderResponse(order), nil
}

// CreateOrderFromCart creates an order from customer's cart.
func (uc *OrderUseCase) CreateOrderFromCart(ctx context.Context, customerID string, req dto.CreateOrderFromCartRequest) (*dto.OrderResponse, error) {
	// Get customer's cart
	cart, err := uc.cartRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, entity.ErrCartNotFound
	}

	if cart.IsEmpty() {
		return nil, entity.ErrEmptyOrder
	}

	// Convert address DTOs to value objects
	shippingAddr, err := valueobject.NewAddress(
		req.ShippingAddress.Street,
		req.ShippingAddress.City,
		req.ShippingAddress.State,
		req.ShippingAddress.PostalCode,
		req.ShippingAddress.Country,
	)
	if err != nil {
		return nil, err
	}

	billingAddr, err := valueobject.NewAddress(
		req.BillingAddress.Street,
		req.BillingAddress.City,
		req.BillingAddress.State,
		req.BillingAddress.PostalCode,
		req.BillingAddress.Country,
	)
	if err != nil {
		return nil, err
	}

	// Generate order ID
	orderID := uc.idGenerator.Generate()

	// Create order from cart using domain service
	order, err := uc.orderService.CreateOrderFromCart(ctx, orderID, cart, shippingAddr, billingAddr)
	if err != nil {
		return nil, err
	}
	order.Notes = req.Notes

	// Clear the cart
	cart.Clear()
	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	return uc.toOrderResponse(order), nil
}

// GetOrder retrieves an order by ID.
func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*dto.OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toOrderResponse(order), nil
}

// GetOrders retrieves all orders with pagination.
func (uc *OrderUseCase) GetOrders(ctx context.Context, limit, offset int) (*dto.OrderListResponse, error) {
	orders, err := uc.orderRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.orderRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		orderResponses[i] = *uc.toOrderResponse(o)
	}

	return &dto.OrderListResponse{
		Orders:  orderResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(orders)) < total,
	}, nil
}

// GetCustomerOrders retrieves orders for a specific customer.
func (uc *OrderUseCase) GetCustomerOrders(ctx context.Context, customerID string, limit, offset int) (*dto.OrderListResponse, error) {
	orders, err := uc.orderRepo.GetByCustomerID(ctx, customerID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.orderRepo.CountByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		orderResponses[i] = *uc.toOrderResponse(o)
	}

	return &dto.OrderListResponse{
		Orders:  orderResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(orders)) < total,
	}, nil
}

// UpdateOrderStatus updates the status of an order.
func (uc *OrderUseCase) UpdateOrderStatus(ctx context.Context, id string, status string) (*dto.OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newStatus := entity.OrderStatus(status)
	if err := order.UpdateStatus(newStatus); err != nil {
		return nil, err
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return uc.toOrderResponse(order), nil
}

// CancelOrder cancels an order and restores stock.
func (uc *OrderUseCase) CancelOrder(ctx context.Context, id string) (*dto.OrderResponse, error) {
	if err := uc.orderService.CancelOrderAndRestoreStock(ctx, id); err != nil {
		return nil, err
	}

	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return uc.toOrderResponse(order), nil
}

// toOrderResponse converts an order entity to a response DTO.
func (uc *OrderUseCase) toOrderResponse(order *entity.Order) *dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice.Amount,
			Currency:    string(item.UnitPrice.Currency),
			TotalPrice:  item.TotalPrice.Amount,
		}
	}

	return &dto.OrderResponse{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Items:       items,
		ItemCount:   order.GetItemCount(),
		TotalAmount: order.TotalAmount.Amount,
		Currency:    string(order.TotalAmount.Currency),
		Status:      string(order.Status),
		ShippingAddress: dto.AddressResponse{
			Street:     order.ShippingAddress.Street,
			City:       order.ShippingAddress.City,
			State:      order.ShippingAddress.State,
			PostalCode: order.ShippingAddress.PostalCode,
			Country:    order.ShippingAddress.Country,
		},
		BillingAddress: dto.AddressResponse{
			Street:     order.BillingAddress.Street,
			City:       order.BillingAddress.City,
			State:      order.BillingAddress.State,
			PostalCode: order.BillingAddress.PostalCode,
			Country:    order.BillingAddress.Country,
		},
		Notes:     order.Notes,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
