package valueobject_test

import (
	"testing"

	"github.com/aruncs31s/esdcshopmodule/internal/domain/valueobject"
)

func TestNewAddress(t *testing.T) {
	tests := []struct {
		name       string
		street     string
		city       string
		state      string
		postalCode string
		country    string
		wantErr    error
	}{
		{
			name:       "valid address",
			street:     "123 Main St",
			city:       "New York",
			state:      "NY",
			postalCode: "10001",
			country:    "USA",
			wantErr:    nil,
		},
		{
			name:       "empty street",
			street:     "",
			city:       "New York",
			state:      "NY",
			postalCode: "10001",
			country:    "USA",
			wantErr:    valueobject.ErrInvalidStreet,
		},
		{
			name:       "empty city",
			street:     "123 Main St",
			city:       "",
			state:      "NY",
			postalCode: "10001",
			country:    "USA",
			wantErr:    valueobject.ErrInvalidCity,
		},
		{
			name:       "empty postal code",
			street:     "123 Main St",
			city:       "New York",
			state:      "NY",
			postalCode: "",
			country:    "USA",
			wantErr:    valueobject.ErrInvalidPostalCode,
		},
		{
			name:       "empty country",
			street:     "123 Main St",
			city:       "New York",
			state:      "NY",
			postalCode: "10001",
			country:    "",
			wantErr:    valueobject.ErrInvalidCountry,
		},
		{
			name:       "empty state is valid",
			street:     "123 Main St",
			city:       "London",
			state:      "",
			postalCode: "SW1A 1AA",
			country:    "UK",
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := valueobject.NewAddress(tt.street, tt.city, tt.state, tt.postalCode, tt.country)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewAddress() unexpected error = %v", err)
				return
			}

			if addr.Street != tt.street {
				t.Errorf("NewAddress() Street = %v, want %v", addr.Street, tt.street)
			}
			if addr.City != tt.city {
				t.Errorf("NewAddress() City = %v, want %v", addr.City, tt.city)
			}
		})
	}
}

func TestAddress_FullAddress(t *testing.T) {
	tests := []struct {
		name     string
		street   string
		city     string
		state    string
		postal   string
		country  string
		expected string
	}{
		{
			name:     "with state",
			street:   "123 Main St",
			city:     "New York",
			state:    "NY",
			postal:   "10001",
			country:  "USA",
			expected: "123 Main St, New York, NY 10001, USA",
		},
		{
			name:     "without state",
			street:   "123 High St",
			city:     "London",
			state:    "",
			postal:   "SW1A 1AA",
			country:  "UK",
			expected: "123 High St, London SW1A 1AA, UK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, _ := valueobject.NewAddress(tt.street, tt.city, tt.state, tt.postal, tt.country)
			result := addr.FullAddress()

			if result != tt.expected {
				t.Errorf("FullAddress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAddress_Equals(t *testing.T) {
	addr1, _ := valueobject.NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	addr2, _ := valueobject.NewAddress("123 Main St", "New York", "NY", "10001", "USA")
	addr3, _ := valueobject.NewAddress("456 Oak Ave", "Los Angeles", "CA", "90001", "USA")

	if !addr1.Equals(addr2) {
		t.Error("Equals() same addresses should be equal")
	}
	if addr1.Equals(addr3) {
		t.Error("Equals() different addresses should not be equal")
	}
}
