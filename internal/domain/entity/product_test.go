package entity_test

import (
	"testing"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/entity"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		productName string
		description string
		price       valueobject.Money
		stock       int
		categoryID  string
		wantErr     error
	}{
		{
			name:        "valid product",
			id:          "prod-1",
			productName: "Test Product",
			description: "A test product",
			price:       valueobject.Money{Amount: 1000, Currency: valueobject.USD},
			stock:       10,
			categoryID:  "cat-1",
			wantErr:     nil,
		},
		{
			name:        "empty id",
			id:          "",
			productName: "Test Product",
			description: "A test product",
			price:       valueobject.Money{Amount: 1000, Currency: valueobject.USD},
			stock:       10,
			categoryID:  "cat-1",
			wantErr:     entity.ErrInvalidProductID,
		},
		{
			name:        "empty name",
			id:          "prod-1",
			productName: "",
			description: "A test product",
			price:       valueobject.Money{Amount: 1000, Currency: valueobject.USD},
			stock:       10,
			categoryID:  "cat-1",
			wantErr:     entity.ErrInvalidProductName,
		},
		{
			name:        "negative stock",
			id:          "prod-1",
			productName: "Test Product",
			description: "A test product",
			price:       valueobject.Money{Amount: 1000, Currency: valueobject.USD},
			stock:       -1,
			categoryID:  "cat-1",
			wantErr:     entity.ErrInvalidStock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := entity.NewProduct(tt.id, tt.productName, tt.description, tt.price, tt.stock, tt.categoryID)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewProduct() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewProduct() unexpected error = %v", err)
				return
			}

			if product.ID != tt.id {
				t.Errorf("NewProduct() ID = %v, want %v", product.ID, tt.id)
			}
			if product.Name != tt.productName {
				t.Errorf("NewProduct() Name = %v, want %v", product.Name, tt.productName)
			}
			if product.Stock != tt.stock {
				t.Errorf("NewProduct() Stock = %v, want %v", product.Stock, tt.stock)
			}
			if !product.Active {
				t.Error("NewProduct() expected Active to be true")
			}
		})
	}
}

func TestProduct_DecreaseStock(t *testing.T) {
	price := valueobject.Money{Amount: 1000, Currency: valueobject.USD}
	product, _ := entity.NewProduct("prod-1", "Test Product", "Description", price, 10, "cat-1")

	tests := []struct {
		name      string
		quantity  int
		wantStock int
		wantErr   error
	}{
		{
			name:      "valid decrease",
			quantity:  3,
			wantStock: 7,
			wantErr:   nil,
		},
		{
			name:      "decrease more than available",
			quantity:  100,
			wantStock: 7, // Stock shouldn't change on error
			wantErr:   entity.ErrInsufficientStock,
		},
		{
			name:      "zero quantity",
			quantity:  0,
			wantStock: 7,
			wantErr:   entity.ErrInvalidQuantity,
		},
		{
			name:      "negative quantity",
			quantity:  -1,
			wantStock: 7,
			wantErr:   entity.ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := product.DecreaseStock(tt.quantity)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("DecreaseStock() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("DecreaseStock() unexpected error = %v", err)
			}

			if product.Stock != tt.wantStock {
				t.Errorf("DecreaseStock() stock = %v, want %v", product.Stock, tt.wantStock)
			}
		})
	}
}

func TestProduct_IncreaseStock(t *testing.T) {
	price := valueobject.Money{Amount: 1000, Currency: valueobject.USD}
	product, _ := entity.NewProduct("prod-1", "Test Product", "Description", price, 10, "cat-1")

	err := product.IncreaseStock(5)
	if err != nil {
		t.Errorf("IncreaseStock() unexpected error = %v", err)
	}
	if product.Stock != 15 {
		t.Errorf("IncreaseStock() stock = %v, want %v", product.Stock, 15)
	}

	err = product.IncreaseStock(0)
	if err != entity.ErrInvalidQuantity {
		t.Errorf("IncreaseStock() with zero error = %v, want %v", err, entity.ErrInvalidQuantity)
	}
}

func TestProduct_IsAvailable(t *testing.T) {
	price := valueobject.Money{Amount: 1000, Currency: valueobject.USD}

	tests := []struct {
		name   string
		active bool
		stock  int
		want   bool
	}{
		{
			name:   "active with stock",
			active: true,
			stock:  10,
			want:   true,
		},
		{
			name:   "active without stock",
			active: true,
			stock:  0,
			want:   false,
		},
		{
			name:   "inactive with stock",
			active: false,
			stock:  10,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, _ := entity.NewProduct("prod-1", "Test", "Desc", price, tt.stock, "cat-1")
			if !tt.active {
				product.Deactivate()
			}

			if got := product.IsAvailable(); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
