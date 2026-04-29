// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package services

import "context"

// InventoryReservation describes the tenant-scoped stock hold the order workflow needs.
type InventoryReservation struct {
	TenantID       int64
	UserID         int64
	OrderID        int64
	TotalCents     int64
	Currency       string
	IdempotencyKey string
}

// InventoryCoordinator is the seam between the order workflow and stock reservation behavior.
type InventoryCoordinator interface {
	Reserve(ctx context.Context, reservation InventoryReservation) error
	// Release must be safe for cleanup retries on duplicate or idempotent order creation flows.
	// Future implementations should treat repeat calls for the same reservation identity as
	// reconciliation, not as a signal to drop the canonical order hold.
	Release(ctx context.Context, reservation InventoryReservation) error
}

// NoopInventoryCoordinator keeps Module 5 focused on the workflow boundary before real stock logic exists.
type NoopInventoryCoordinator struct{}

// NewNoopInventoryCoordinator creates an inventory stub that always succeeds.
func NewNoopInventoryCoordinator() NoopInventoryCoordinator {
	return NoopInventoryCoordinator{}
}

// Reserve simulates successfully acquiring stock for an order.
func (NoopInventoryCoordinator) Reserve(context.Context, InventoryReservation) error {
	return nil
}

// Release simulates successfully dropping a hold on stock after an order fails or cancels.
func (NoopInventoryCoordinator) Release(context.Context, InventoryReservation) error {
	return nil
}
