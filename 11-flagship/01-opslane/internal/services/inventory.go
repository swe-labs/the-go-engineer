// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package services

import "context"

// InventoryReservation (Struct): tenant-scoped stock hold description for the order workflow
type InventoryReservation struct {
	TenantID       int64
	UserID         int64
	OrderID        int64
	TotalCents     int64
	Currency       string
	IdempotencyKey string
}

// InventoryCoordinator (Interface): seam between the order workflow and stock reservation behavior
type InventoryCoordinator interface {
	Reserve(ctx context.Context, reservation InventoryReservation) error
	// Release must be safe for cleanup retries on duplicate or idempotent order creation flows.
	// Future implementations should treat repeat calls for the same reservation identity as
	// reconciliation, not as a signal to drop the canonical order hold.
	Release(ctx context.Context, reservation InventoryReservation) error
}

// NoopInventoryCoordinator (Struct): stub inventory coordinator for development before real stock logic exists
type NoopInventoryCoordinator struct{}

// NewNoopInventoryCoordinator (Constructor): creates an inventory stub that always succeeds
func NewNoopInventoryCoordinator() NoopInventoryCoordinator {
	return NoopInventoryCoordinator{}
}

// Reserve (Method): no-op that simulates successfully acquiring stock
func (NoopInventoryCoordinator) Reserve(context.Context, InventoryReservation) error {
	return nil
}

// Release (Method): no-op that simulates successfully releasing stock hold
func (NoopInventoryCoordinator) Release(context.Context, InventoryReservation) error {
	return nil
}
