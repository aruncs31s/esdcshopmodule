package service

import (
	"context"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/repository"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// OrderService provides domain-level operations for orders.
// This follows the Single Responsibility Principle (SRP) from SOLID.
type OrderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	cartRepo    repository.CartRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	cartRepo repository.CartRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
	}
}

// CreateOrderFromCart creates an order from a customer's cart.
func (s *OrderService) CreateOrderFromCart(
	ctx context.Context,
	orderID string,
	cart *entity.Cart,
	shippingAddr, billingAddr valueobject.Address,
) (*entity.Order, error) {
	if cart.IsEmpty() {
		return nil, entity.ErrEmptyOrder
	}

	// Convert cart items to order items
	orderItems := make([]entity.OrderItem, 0, len(cart.Items))
	for _, cartItem := range cart.Items {
		product, err := s.productRepo.GetByID(ctx, cartItem.ProductID)
		if err != nil {
			return nil, err
		}

		// Check stock availability
		if product.Stock < cartItem.Quantity {
			return nil, entity.ErrInsufficientStock
		}

		orderItem := entity.OrderItem{
			ProductID:   cartItem.ProductID,
			ProductName: product.Name,
			Quantity:    cartItem.Quantity,
			UnitPrice:   cartItem.UnitPrice,
			TotalPrice:  cartItem.TotalPrice,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create the order
	order, err := entity.NewOrder(orderID, cart.CustomerID, orderItems, shippingAddr, billingAddr)
	if err != nil {
		return nil, err
	}

	// Decrease stock for each product
	for _, item := range orderItems {
		product, _ := s.productRepo.GetByID(ctx, item.ProductID)
		if err := product.DecreaseStock(item.Quantity); err != nil {
			return nil, err
		}
		if err := s.productRepo.Update(ctx, product); err != nil {
			return nil, err
		}
	}

	// Save the order
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// CancelOrderAndRestoreStock cancels an order and restores the product stock.
func (s *OrderService) CancelOrderAndRestoreStock(ctx context.Context, orderID string) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !order.IsCancellable() {
		return entity.ErrInvalidOrderStatus
	}

	// Restore stock for each product
	for _, item := range order.Items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			continue // Product might have been deleted
		}
		if err := product.IncreaseStock(item.Quantity); err != nil {
			return err
		}
		if err := s.productRepo.Update(ctx, product); err != nil {
			return err
		}
	}

	// Cancel the order
	if err := order.Cancel(); err != nil {
		return err
	}

	return s.orderRepo.Update(ctx, order)
}

// CalculateOrderTotal calculates the total amount for a list of order items.
func (s *OrderService) CalculateOrderTotal(items []entity.OrderItem) (valueobject.Money, error) {
	if len(items) == 0 {
		return valueobject.Money{}, entity.ErrEmptyOrder
	}

	total := valueobject.Money{Amount: 0, Currency: items[0].TotalPrice.Currency}
	for _, item := range items {
		var err error
		total, err = total.Add(item.TotalPrice)
		if err != nil {
			return valueobject.Money{}, err
		}
	}
	return total, nil
}
