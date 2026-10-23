package main

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidOrder       = errors.New("order must contain a customer and at least one item")
	ErrInvalidItem        = errors.New("item name cannot be empty and amount must be positive")
	ErrOrderNotFound      = errors.New("order not found")
	ErrPaymentFailed      = errors.New("payment processing rejected")
	ErrOrderAlreadyExists = errors.New("order with this id already exists")
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusCancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	SKU        string `json:"sku"`
	Quantity   int    `json:"quantity"`
	PriceCents int    `json:"price_cents"`
}

type Order struct {
	ID         string      `json:"id"`
	Customer   string      `json:"customer"`
	Email      string      `json:"email"`
	Items      []OrderItem `json:"items"`
	Subtotal   int         `json:"subtotal_cents"`
	Discount   int         `json:"discount_cents"`
	Total      int         `json:"total_cents"`
	Status     OrderStatus `json:"status"`
	PaymentRef string      `json:"payment_ref,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}

// CalculateTotals calculates the subtotal, applies discount rules (10% off for orders >= $100),
// and validates order invariants.
func (o *Order) CalculateTotals() error {
	// TODO: Task 1 - Implement domain calculation and invariants:
	// - Customer and Email must not be empty (ErrInvalidOrder)
	// - Must contain at least one item (ErrInvalidOrder)
	// - Items SKU cannot be empty, Quantity > 0, PriceCents > 0 (ErrInvalidItem)
	// - 10% discount on Subtotal >= 10,000 cents ($100)
	if strings.TrimSpace(o.Customer) == "" || strings.TrimSpace(o.Email) == "" || len(o.Items) == 0 {
		return ErrInvalidOrder
	}
	return nil
}
