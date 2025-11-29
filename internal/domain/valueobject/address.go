package valueobject

import "errors"

// Address represents a shipping or billing address.
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

var (
	ErrInvalidStreet     = errors.New("street cannot be empty")
	ErrInvalidCity       = errors.New("city cannot be empty")
	ErrInvalidPostalCode = errors.New("postal code cannot be empty")
	ErrInvalidCountry    = errors.New("country cannot be empty")
)

// NewAddress creates a new Address value object with validation.
func NewAddress(street, city, state, postalCode, country string) (Address, error) {
	if street == "" {
		return Address{}, ErrInvalidStreet
	}
	if city == "" {
		return Address{}, ErrInvalidCity
	}
	if postalCode == "" {
		return Address{}, ErrInvalidPostalCode
	}
	if country == "" {
		return Address{}, ErrInvalidCountry
	}

	return Address{
		Street:     street,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
	}, nil
}

// FullAddress returns the complete address as a formatted string.
func (a Address) FullAddress() string {
	addr := a.Street + ", " + a.City
	if a.State != "" {
		addr += ", " + a.State
	}
	addr += " " + a.PostalCode + ", " + a.Country
	return addr
}

// Equals checks if two addresses are equal.
func (a Address) Equals(other Address) bool {
	return a.Street == other.Street &&
		a.City == other.City &&
		a.State == other.State &&
		a.PostalCode == other.PostalCode &&
		a.Country == other.Country
}
