// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

// PaymentStatus represents the lifecycle state of a payment attempt.
type PaymentStatus string

const (
	// PaymentStatusPending means the payment intent is recorded but not yet sent to the gateway.
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusAuthorized means the gateway placed a hold on the funds.
	PaymentStatusAuthorized PaymentStatus = "authorized"
	// PaymentStatusSettled means the funds have been successfully captured.
	PaymentStatusSettled PaymentStatus = "settled"
	// PaymentStatusFailed means the gateway declined the transaction.
	PaymentStatusFailed PaymentStatus = "failed"
	// PaymentStatusRefunded means the payment was returned to the customer.
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Payment records a financial transaction attempt against an order.
// Orders may have multiple payments if previous attempts failed.
type Payment struct {
	ID                int64         `json:"id"`
	TenantID          int64         `json:"tenant_id"`
	OrderID           int64         `json:"order_id"`
	Status            PaymentStatus `json:"status"`
	ProviderReference string        `json:"provider_reference"`
	AmountCents       int64         `json:"amount_cents"`
	FailureReason     string        `json:"failure_reason,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}
