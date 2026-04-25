// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusAuthorized PaymentStatus = "authorized"
	PaymentStatusSettled    PaymentStatus = "settled"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

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
