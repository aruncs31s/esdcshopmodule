package valueobject

import (
	"errors"
	"fmt"
)

// Currency represents a currency code.
type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	INR Currency = "INR"
)

// Money represents a monetary value with currency.
type Money struct {
	Amount   int64    // Amount in smallest unit (cents, paise, etc.)
	Currency Currency
}

var (
	ErrInvalidAmount   = errors.New("amount cannot be negative")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrCurrencyMismatch = errors.New("currency mismatch")
)

// NewMoney creates a new Money value object.
func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, ErrInvalidAmount
	}
	if !isValidCurrency(currency) {
		return Money{}, ErrInvalidCurrency
	}
	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

// Add adds two Money values.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

// Subtract subtracts two Money values.
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	if m.Amount < other.Amount {
		return Money{}, ErrInvalidAmount
	}
	return Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}, nil
}

// Multiply multiplies the money by a quantity.
func (m Money) Multiply(quantity int) Money {
	return Money{
		Amount:   m.Amount * int64(quantity),
		Currency: m.Currency,
	}
}

// String returns a string representation of the money.
func (m Money) String() string {
	return fmt.Sprintf("%s %.2f", m.Currency, float64(m.Amount)/100)
}

// Equals checks if two Money values are equal.
func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

// IsZero checks if the money amount is zero.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

func isValidCurrency(currency Currency) bool {
	switch currency {
	case USD, EUR, GBP, INR:
		return true
	default:
		return false
	}
}
