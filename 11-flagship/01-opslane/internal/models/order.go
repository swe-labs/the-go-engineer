// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package models

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusFailed     OrderStatus = "failed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

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
