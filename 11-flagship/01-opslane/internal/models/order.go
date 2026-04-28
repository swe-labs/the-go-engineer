// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

// OrderStatus defines the state machine states for a customer order.
type OrderStatus string

const (
	// OrderStatusPending means the order is created but not yet paid or processed.
	OrderStatusPending OrderStatus = "pending"
	// OrderStatusProcessing means the payment is actively being handled by the gateway.
	OrderStatusProcessing OrderStatus = "processing"
	// OrderStatusPaid means the payment succeeded and the order is ready for fulfillment.
	OrderStatusPaid OrderStatus = "paid"
	// OrderStatusFailed means the payment failed and the order cannot proceed.
	OrderStatusFailed OrderStatus = "failed"
	// OrderStatusCancelled means the order was aborted before or after payment.
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a customer's purchase intent. It acts as the aggregate root
// for Payments and Inventory reservations.
type Order struct {
	ID             int64       `json:"id"`
	TenantID       int64       `json:"tenant_id"`
	UserID         int64       `json:"user_id"`
	Status         OrderStatus `json:"status"`
	TotalCents     int64       `json:"total_cents"`
	Currency       string      `json:"currency"`
	IdempotencyKey string      `json:"idempotency_key"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
