// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package events

import (
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

// Type represents the routing key for a specific business event.
// Subscribers use these keys to bind handlers to specific lifecycle moments.
type Type string

const (
	// TypeOrderCreated is emitted immediately after a new order is persisted.
	TypeOrderCreated Type = "order.created"
	// TypeOrderStatusChanged is emitted whenever an order transitions (e.g. Pending -> Paid).
	TypeOrderStatusChanged Type = "order.status_changed"
	// TypePaymentRequested is emitted when a payment needs to be authorized by the gateway.
	TypePaymentRequested Type = "payment.requested"
	// TypePaymentSettled is emitted when the gateway confirms funds are captured.
	TypePaymentSettled Type = "payment.settled"
	// TypePaymentFailed is emitted when a payment attempt is rejected by the gateway.
	TypePaymentFailed Type = "payment.failed"
	// TypeNotificationRequested is emitted to trigger an asynchronous customer communication.
	TypeNotificationRequested Type = "notification.requested"
)

// Event is the small, explicit message that crosses the async boundary.
type Event struct {
	ID                string
	Type              Type
	TenantID          int64
	UserID            int64
	OrderID           int64
	PaymentID         int64
	OrderStatus       models.OrderStatus
	PaymentStatus     models.PaymentStatus
	ProviderReference string
	AmountCents       int64
	OccurredAt        time.Time
	Metadata          map[string]string
}

// WithOccurredAt returns a copy of the event with the timestamp populated.
// If the event already has a timestamp, it is preserved.
func (e Event) WithOccurredAt(now time.Time) Event {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = now
	}

	return e
}
