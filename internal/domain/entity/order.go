package entity

import (
	"time"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// OrderItem represents an item in an order.
type OrderItem struct {
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   valueobject.Money
	TotalPrice  valueobject.Money
}

// Order represents an order entity in the shop domain.
type Order struct {
	ID              string
	CustomerID      string
	Items           []OrderItem
	TotalAmount     valueobject.Money
	Status          OrderStatus
	ShippingAddress valueobject.Address
	BillingAddress  valueobject.Address
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewOrder creates a new Order with validation.
func NewOrder(id, customerID string, items []OrderItem, shippingAddr, billingAddr valueobject.Address) (*Order, error) {
	if id == "" {
		return nil, ErrInvalidOrderID
	}
	if customerID == "" {
		return nil, ErrInvalidCustomerID
	}
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}

	now := time.Now()
	order := &Order{
		ID:              id,
		CustomerID:      customerID,
		Items:           items,
		Status:          OrderStatusPending,
		ShippingAddress: shippingAddr,
		BillingAddress:  billingAddr,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := order.calculateTotal(); err != nil {
		return nil, err
	}

	return order, nil
}

// calculateTotal calculates the total order amount.
func (o *Order) calculateTotal() error {
	if len(o.Items) == 0 {
		return ErrEmptyOrder
	}

	total := valueobject.Money{Amount: 0, Currency: o.Items[0].TotalPrice.Currency}
	for _, item := range o.Items {
		var err error
		total, err = total.Add(item.TotalPrice)
		if err != nil {
			return err
		}
	}
	o.TotalAmount = total
	return nil
}

// CanTransitionTo checks if the order can transition to a new status.
func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending:    {OrderStatusConfirmed, OrderStatusCancelled},
		OrderStatusConfirmed:  {OrderStatusProcessing, OrderStatusCancelled},
		OrderStatusProcessing: {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:    {OrderStatusDelivered},
		OrderStatusDelivered:  {OrderStatusRefunded},
		OrderStatusCancelled:  {},
		OrderStatusRefunded:   {},
	}

	allowedStatuses, exists := validTransitions[o.Status]
	if !exists {
		return false
	}

	for _, s := range allowedStatuses {
		if s == newStatus {
			return true
		}
	}
	return false
}

// UpdateStatus updates the order status.
func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	if !o.CanTransitionTo(newStatus) {
		return ErrInvalidOrderStatus
	}
	o.Status = newStatus
	o.UpdatedAt = time.Now()
	return nil
}

// Confirm confirms the order.
func (o *Order) Confirm() error {
	return o.UpdateStatus(OrderStatusConfirmed)
}

// Process marks the order as processing.
func (o *Order) Process() error {
	return o.UpdateStatus(OrderStatusProcessing)
}

// Ship marks the order as shipped.
func (o *Order) Ship() error {
	return o.UpdateStatus(OrderStatusShipped)
}

// Deliver marks the order as delivered.
func (o *Order) Deliver() error {
	return o.UpdateStatus(OrderStatusDelivered)
}

// Cancel cancels the order.
func (o *Order) Cancel() error {
	return o.UpdateStatus(OrderStatusCancelled)
}

// Refund refunds the order.
func (o *Order) Refund() error {
	return o.UpdateStatus(OrderStatusRefunded)
}

// IsCancellable checks if the order can be cancelled.
func (o *Order) IsCancellable() bool {
	return o.CanTransitionTo(OrderStatusCancelled)
}

// GetItemCount returns the total number of items in the order.
func (o *Order) GetItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}
