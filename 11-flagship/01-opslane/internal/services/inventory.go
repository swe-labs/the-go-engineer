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
	Release(ctx context.Context, reservation InventoryReservation) error
}

// NoopInventoryCoordinator keeps Module 5 focused on the workflow boundary before real stock logic exists.
type NoopInventoryCoordinator struct{}

func NewNoopInventoryCoordinator() NoopInventoryCoordinator {
	return NoopInventoryCoordinator{}
}

func (NoopInventoryCoordinator) Reserve(context.Context, InventoryReservation) error {
	return nil
}

func (NoopInventoryCoordinator) Release(context.Context, InventoryReservation) error {
	return nil
}
