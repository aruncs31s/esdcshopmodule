package entity_test

import (
	"testing"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

func TestNewCart(t *testing.T) {
	tests := []struct {
		name       string
		cartID     string
		customerID string
		currency   valueobject.Currency
		wantErr    error
	}{
		{
			name:       "valid cart",
			cartID:     "cart-1",
			customerID: "cust-1",
			currency:   valueobject.USD,
			wantErr:    nil,
		},
		{
			name:       "empty cart id",
			cartID:     "",
			customerID: "cust-1",
			currency:   valueobject.USD,
			wantErr:    entity.ErrInvalidCartID,
		},
		{
			name:       "empty customer id",
			cartID:     "cart-1",
			customerID: "",
			currency:   valueobject.USD,
			wantErr:    entity.ErrInvalidCustomerID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart, err := entity.NewCart(tt.cartID, tt.customerID, tt.currency)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewCart() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewCart() unexpected error = %v", err)
				return
			}

			if cart.ID != tt.cartID {
				t.Errorf("NewCart() ID = %v, want %v", cart.ID, tt.cartID)
			}
			if cart.CustomerID != tt.customerID {
				t.Errorf("NewCart() CustomerID = %v, want %v", cart.CustomerID, tt.customerID)
			}
			if !cart.IsEmpty() {
				t.Error("NewCart() cart should be empty")
			}
		})
	}
}

func TestCart_AddItem(t *testing.T) {
	cart, _ := entity.NewCart("cart-1", "cust-1", valueobject.USD)
	price := valueobject.Money{Amount: 1000, Currency: valueobject.USD}

	// Add first item
	err := cart.AddItem("prod-1", 2, price)
	if err != nil {
		t.Errorf("AddItem() error = %v", err)
	}

	if cart.GetItemCount() != 2 {
		t.Errorf("GetItemCount() = %v, want %v", cart.GetItemCount(), 2)
	}

	if cart.Total.Amount != 2000 {
		t.Errorf("Total.Amount = %v, want %v", cart.Total.Amount, 2000)
	}

	// Add same item again (should increase quantity)
	err = cart.AddItem("prod-1", 3, price)
	if err != nil {
		t.Errorf("AddItem() same item error = %v", err)
	}

	if cart.GetItemCount() != 5 {
		t.Errorf("GetItemCount() after adding same item = %v, want %v", cart.GetItemCount(), 5)
	}

	if cart.Total.Amount != 5000 {
		t.Errorf("Total.Amount after adding same item = %v, want %v", cart.Total.Amount, 5000)
	}

	// Add different item
	price2 := valueobject.Money{Amount: 500, Currency: valueobject.USD}
	err = cart.AddItem("prod-2", 2, price2)
	if err != nil {
		t.Errorf("AddItem() different item error = %v", err)
	}

	if cart.GetItemCount() != 7 {
		t.Errorf("GetItemCount() after adding different item = %v, want %v", cart.GetItemCount(), 7)
	}

	if cart.Total.Amount != 6000 {
		t.Errorf("Total.Amount after adding different item = %v, want %v", cart.Total.Amount, 6000)
	}
}

func TestCart_RemoveItem(t *testing.T) {
	cart, _ := entity.NewCart("cart-1", "cust-1", valueobject.USD)
	price := valueobject.Money{Amount: 1000, Currency: valueobject.USD}

	cart.AddItem("prod-1", 2, price)
	cart.AddItem("prod-2", 3, price)

	// Remove item
	err := cart.RemoveItem("prod-1")
	if err != nil {
		t.Errorf("RemoveItem() error = %v", err)
	}

	if cart.GetItemCount() != 3 {
		t.Errorf("GetItemCount() after remove = %v, want %v", cart.GetItemCount(), 3)
	}

	// Try to remove non-existent item
	err = cart.RemoveItem("prod-999")
	if err != entity.ErrProductNotFound {
		t.Errorf("RemoveItem() non-existent error = %v, want %v", err, entity.ErrProductNotFound)
	}
}

func TestCart_UpdateItemQuantity(t *testing.T) {
	cart, _ := entity.NewCart("cart-1", "cust-1", valueobject.USD)
	price := valueobject.Money{Amount: 1000, Currency: valueobject.USD}

	cart.AddItem("prod-1", 2, price)

	// Update quantity
	err := cart.UpdateItemQuantity("prod-1", 5)
	if err != nil {
		t.Errorf("UpdateItemQuantity() error = %v", err)
	}

	if cart.GetItemCount() != 5 {
		t.Errorf("GetItemCount() after update = %v, want %v", cart.GetItemCount(), 5)
	}

	if cart.Total.Amount != 5000 {
		t.Errorf("Total.Amount after update = %v, want %v", cart.Total.Amount, 5000)
	}

	// Update to zero should remove item
	err = cart.UpdateItemQuantity("prod-1", 0)
	if err != nil {
		t.Errorf("UpdateItemQuantity() to zero error = %v", err)
	}

	if !cart.IsEmpty() {
		t.Error("Cart should be empty after updating to zero quantity")
	}
}

func TestCart_Clear(t *testing.T) {
	cart, _ := entity.NewCart("cart-1", "cust-1", valueobject.USD)
	price := valueobject.Money{Amount: 1000, Currency: valueobject.USD}

	cart.AddItem("prod-1", 2, price)
	cart.AddItem("prod-2", 3, price)

	cart.Clear()

	if !cart.IsEmpty() {
		t.Error("Cart should be empty after Clear()")
	}

	if cart.Total.Amount != 0 {
		t.Errorf("Total.Amount after Clear() = %v, want %v", cart.Total.Amount, 0)
	}
}
