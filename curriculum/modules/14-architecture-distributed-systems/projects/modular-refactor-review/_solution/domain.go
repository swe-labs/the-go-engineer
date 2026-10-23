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
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
	PriceCents  int    `json:"price_cents"`
}

type Order struct {
	ID          string      `json:"id"`
	Customer    string      `json:"customer"`
	Email       string      `json:"email"`
	Items       []OrderItem `json:"items"`
	Subtotal    int         `json:"subtotal_cents"`
	Discount    int         `json:"discount_cents"`
	Total       int         `json:"total_cents"`
	Status      OrderStatus `json:"status"`
	PaymentRef  string      `json:"payment_ref,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
}

// CalculateTotals calculates the subtotal, applies discount rules (10% off for orders >= $100),
// and validates order invariants.
func (o *Order) CalculateTotals() error {
	if strings.TrimSpace(o.Customer) == "" || strings.TrimSpace(o.Email) == "" || len(o.Items) == 0 {
		return ErrInvalidOrder
	}

	subtotal := 0
	for _, item := range o.Items {
		if strings.TrimSpace(item.SKU) == "" || item.Quantity <= 0 || item.PriceCents <= 0 {
			return ErrInvalidItem
		}
		subtotal += item.PriceCents * item.Quantity
	}

	o.Subtotal = subtotal

	// Business Rule: 10% discount for orders of $100 (10,000 cents) or more
	discount := 0
	if subtotal >= 10000 {
		discount = subtotal / 10
	}
	o.Discount = discount
	o.Total = subtotal - discount
	return nil
}
