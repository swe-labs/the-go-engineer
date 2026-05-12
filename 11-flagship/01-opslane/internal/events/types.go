// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package events

import (
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/models"
)

// Type (Type): routing key for a specific business event; used by subscribers to bind handlers
type Type string

const (
	// TypeOrderCreated (Constant): emitted immediately after a new order is persisted
	TypeOrderCreated Type = "order.created"
	// TypeOrderStatusChanged (Constant): emitted when an order transitions status
	TypeOrderStatusChanged Type = "order.status_changed"
	// TypePaymentRequested (Constant): emitted when a payment needs gateway authorization
	TypePaymentRequested Type = "payment.requested"
	// TypePaymentSettled (Constant): emitted when the gateway confirms funds are captured
	TypePaymentSettled Type = "payment.settled"
	// TypePaymentFailed (Constant): emitted when a payment attempt is rejected by the gateway
	TypePaymentFailed Type = "payment.failed"
	// TypeNotificationRequested (Constant): emitted to trigger asynchronous customer communication
	TypeNotificationRequested Type = "notification.requested"
)

// Event (Struct): small, explicit message that crosses the async event boundary
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

// WithOccurredAt (Method): returns a copy of the event with the timestamp populated if absent
func (e Event) WithOccurredAt(now time.Time) Event {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = now
	}

	return e
}
