package valueobject_test

import (
	"testing"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency valueobject.Currency
		wantErr  error
	}{
		{
			name:     "valid money",
			amount:   1000,
			currency: valueobject.USD,
			wantErr:  nil,
		},
		{
			name:     "zero amount",
			amount:   0,
			currency: valueobject.USD,
			wantErr:  nil,
		},
		{
			name:     "negative amount",
			amount:   -100,
			currency: valueobject.USD,
			wantErr:  valueobject.ErrInvalidAmount,
		},
		{
			name:     "invalid currency",
			amount:   1000,
			currency: "INVALID",
			wantErr:  valueobject.ErrInvalidCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := valueobject.NewMoney(tt.amount, tt.currency)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewMoney() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewMoney() unexpected error = %v", err)
				return
			}

			if money.Amount != tt.amount {
				t.Errorf("NewMoney() Amount = %v, want %v", money.Amount, tt.amount)
			}
			if money.Currency != tt.currency {
				t.Errorf("NewMoney() Currency = %v, want %v", money.Currency, tt.currency)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	m1 := valueobject.Money{Amount: 1000, Currency: valueobject.USD}
	m2 := valueobject.Money{Amount: 500, Currency: valueobject.USD}
	m3 := valueobject.Money{Amount: 500, Currency: valueobject.EUR}

	// Add same currency
	result, err := m1.Add(m2)
	if err != nil {
		t.Errorf("Add() error = %v", err)
	}
	if result.Amount != 1500 {
		t.Errorf("Add() Amount = %v, want %v", result.Amount, 1500)
	}

	// Add different currency
	_, err = m1.Add(m3)
	if err != valueobject.ErrCurrencyMismatch {
		t.Errorf("Add() different currency error = %v, want %v", err, valueobject.ErrCurrencyMismatch)
	}
}

func TestMoney_Subtract(t *testing.T) {
	m1 := valueobject.Money{Amount: 1000, Currency: valueobject.USD}
	m2 := valueobject.Money{Amount: 500, Currency: valueobject.USD}
	m3 := valueobject.Money{Amount: 2000, Currency: valueobject.USD}
	m4 := valueobject.Money{Amount: 500, Currency: valueobject.EUR}

	// Subtract valid amount
	result, err := m1.Subtract(m2)
	if err != nil {
		t.Errorf("Subtract() error = %v", err)
	}
	if result.Amount != 500 {
		t.Errorf("Subtract() Amount = %v, want %v", result.Amount, 500)
	}

	// Subtract more than available
	_, err = m1.Subtract(m3)
	if err != valueobject.ErrInvalidAmount {
		t.Errorf("Subtract() more than available error = %v, want %v", err, valueobject.ErrInvalidAmount)
	}

	// Subtract different currency
	_, err = m1.Subtract(m4)
	if err != valueobject.ErrCurrencyMismatch {
		t.Errorf("Subtract() different currency error = %v, want %v", err, valueobject.ErrCurrencyMismatch)
	}
}

func TestMoney_Multiply(t *testing.T) {
	m := valueobject.Money{Amount: 1000, Currency: valueobject.USD}

	result := m.Multiply(3)
	if result.Amount != 3000 {
		t.Errorf("Multiply() Amount = %v, want %v", result.Amount, 3000)
	}
	if result.Currency != valueobject.USD {
		t.Errorf("Multiply() Currency = %v, want %v", result.Currency, valueobject.USD)
	}
}

func TestMoney_Equals(t *testing.T) {
	m1 := valueobject.Money{Amount: 1000, Currency: valueobject.USD}
	m2 := valueobject.Money{Amount: 1000, Currency: valueobject.USD}
	m3 := valueobject.Money{Amount: 500, Currency: valueobject.USD}
	m4 := valueobject.Money{Amount: 1000, Currency: valueobject.EUR}

	if !m1.Equals(m2) {
		t.Error("Equals() same values should be equal")
	}
	if m1.Equals(m3) {
		t.Error("Equals() different amounts should not be equal")
	}
	if m1.Equals(m4) {
		t.Error("Equals() different currencies should not be equal")
	}
}

func TestMoney_IsZero(t *testing.T) {
	zero := valueobject.Money{Amount: 0, Currency: valueobject.USD}
	nonZero := valueobject.Money{Amount: 1000, Currency: valueobject.USD}

	if !zero.IsZero() {
		t.Error("IsZero() zero amount should return true")
	}
	if nonZero.IsZero() {
		t.Error("IsZero() non-zero amount should return false")
	}
}

func TestMoney_String(t *testing.T) {
	m := valueobject.Money{Amount: 1050, Currency: valueobject.USD}

	str := m.String()
	if str != "USD 10.50" {
		t.Errorf("String() = %v, want %v", str, "USD 10.50")
	}
}
