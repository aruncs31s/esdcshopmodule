package entity_test

import (
	"testing"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

func TestNewOrder(t *testing.T) {
	shippingAddr, _ := valueobject.NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	billingAddr, _ := valueobject.NewAddress("456 Oak Ave", "Los Angeles", "CA", "90001", "USA")

	items := []entity.OrderItem{
		{
			ProductID:   "prod-1",
			ProductName: "Product 1",
			Quantity:    2,
			UnitPrice:   valueobject.Money{Amount: 1000, Currency: valueobject.USD},
			TotalPrice:  valueobject.Money{Amount: 2000, Currency: valueobject.USD},
		},
	}

	tests := []struct {
		name       string
		orderID    string
		customerID string
		items      []entity.OrderItem
		wantErr    error
	}{
		{
			name:       "valid order",
			orderID:    "order-1",
			customerID: "cust-1",
			items:      items,
			wantErr:    nil,
		},
		{
			name:       "empty order id",
			orderID:    "",
			customerID: "cust-1",
			items:      items,
			wantErr:    entity.ErrInvalidOrderID,
		},
		{
			name:       "empty customer id",
			orderID:    "order-1",
			customerID: "",
			items:      items,
			wantErr:    entity.ErrInvalidCustomerID,
		},
		{
			name:       "empty items",
			orderID:    "order-1",
			customerID: "cust-1",
			items:      []entity.OrderItem{},
			wantErr:    entity.ErrEmptyOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := entity.NewOrder(tt.orderID, tt.customerID, tt.items, shippingAddr, billingAddr)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewOrder() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewOrder() unexpected error = %v", err)
				return
			}

			if order.ID != tt.orderID {
				t.Errorf("NewOrder() ID = %v, want %v", order.ID, tt.orderID)
			}
			if order.CustomerID != tt.customerID {
				t.Errorf("NewOrder() CustomerID = %v, want %v", order.CustomerID, tt.customerID)
			}
			if order.Status != entity.OrderStatusPending {
				t.Errorf("NewOrder() Status = %v, want %v", order.Status, entity.OrderStatusPending)
			}
		})
	}
}

func TestOrder_StatusTransitions(t *testing.T) {
	shippingAddr, _ := valueobject.NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	billingAddr, _ := valueobject.NewAddress("456 Oak Ave", "Los Angeles", "CA", "90001", "USA")

	items := []entity.OrderItem{
		{
			ProductID:   "prod-1",
			ProductName: "Product 1",
			Quantity:    2,
			UnitPrice:   valueobject.Money{Amount: 1000, Currency: valueobject.USD},
			TotalPrice:  valueobject.Money{Amount: 2000, Currency: valueobject.USD},
		},
	}

	t.Run("valid status flow", func(t *testing.T) {
		order, _ := entity.NewOrder("order-1", "cust-1", items, shippingAddr, billingAddr)

		// Pending -> Confirmed
		if err := order.Confirm(); err != nil {
			t.Errorf("Confirm() error = %v", err)
		}
		if order.Status != entity.OrderStatusConfirmed {
			t.Errorf("Status after Confirm() = %v, want %v", order.Status, entity.OrderStatusConfirmed)
		}

		// Confirmed -> Processing
		if err := order.Process(); err != nil {
			t.Errorf("Process() error = %v", err)
		}
		if order.Status != entity.OrderStatusProcessing {
			t.Errorf("Status after Process() = %v, want %v", order.Status, entity.OrderStatusProcessing)
		}

		// Processing -> Shipped
		if err := order.Ship(); err != nil {
			t.Errorf("Ship() error = %v", err)
		}
		if order.Status != entity.OrderStatusShipped {
			t.Errorf("Status after Ship() = %v, want %v", order.Status, entity.OrderStatusShipped)
		}

		// Shipped -> Delivered
		if err := order.Deliver(); err != nil {
			t.Errorf("Deliver() error = %v", err)
		}
		if order.Status != entity.OrderStatusDelivered {
			t.Errorf("Status after Deliver() = %v, want %v", order.Status, entity.OrderStatusDelivered)
		}
	})

	t.Run("invalid status transition", func(t *testing.T) {
		order, _ := entity.NewOrder("order-2", "cust-1", items, shippingAddr, billingAddr)

		// Pending -> Shipped (invalid)
		if err := order.Ship(); err != entity.ErrInvalidOrderStatus {
			t.Errorf("Ship() from Pending error = %v, want %v", err, entity.ErrInvalidOrderStatus)
		}
	})

	t.Run("cancel order", func(t *testing.T) {
		order, _ := entity.NewOrder("order-3", "cust-1", items, shippingAddr, billingAddr)

		if err := order.Cancel(); err != nil {
			t.Errorf("Cancel() error = %v", err)
		}
		if order.Status != entity.OrderStatusCancelled {
			t.Errorf("Status after Cancel() = %v, want %v", order.Status, entity.OrderStatusCancelled)
		}

		// Cannot cancel already cancelled order
		if order.IsCancellable() {
			t.Error("IsCancellable() = true for cancelled order")
		}
	})
}

func TestOrder_GetItemCount(t *testing.T) {
	shippingAddr, _ := valueobject.NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	billingAddr, _ := valueobject.NewAddress("456 Oak Ave", "Los Angeles", "CA", "90001", "USA")

	items := []entity.OrderItem{
		{
			ProductID:   "prod-1",
			ProductName: "Product 1",
			Quantity:    2,
			UnitPrice:   valueobject.Money{Amount: 1000, Currency: valueobject.USD},
			TotalPrice:  valueobject.Money{Amount: 2000, Currency: valueobject.USD},
		},
		{
			ProductID:   "prod-2",
			ProductName: "Product 2",
			Quantity:    3,
			UnitPrice:   valueobject.Money{Amount: 500, Currency: valueobject.USD},
			TotalPrice:  valueobject.Money{Amount: 1500, Currency: valueobject.USD},
		},
	}

	order, _ := entity.NewOrder("order-1", "cust-1", items, shippingAddr, billingAddr)

	if got := order.GetItemCount(); got != 5 {
		t.Errorf("GetItemCount() = %v, want %v", got, 5)
	}
}
