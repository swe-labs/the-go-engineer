// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

// OrderRepository is the minimum persistence surface the order workflow needs.
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, tenantID, orderID int64) (models.Order, error)
	GetOrderByIdempotencyKey(ctx context.Context, tenantID int64, idempotencyKey string) (models.Order, error)
	UpdateOrderStatus(ctx context.Context, tenantID, orderID int64, status models.OrderStatus) (models.Order, error)
}

// CreateOrderInput carries the tenant-scoped values required to begin an order workflow.
type CreateOrderInput struct {
	TenantID       int64
	UserID         int64
	TotalCents     int64
	Currency       string
	IdempotencyKey string
}

// CreateOrderResult differentiates a brand-new order from an idempotent replay.
type CreateOrderResult struct {
	Order   models.Order
	Created bool
}

// TransitionOrderRequest describes a tenant-scoped order status change.
type TransitionOrderRequest struct {
	TenantID int64
	OrderID  int64
	Status   models.OrderStatus
}

// OrderService moves order creation and state changes out of handlers and into business workflow code.
type OrderService struct {
	orders    OrderRepository
	inventory InventoryCoordinator
}

// NewOrderService initializes the core business logic for order management.
// If no inventory coordinator is provided, it defaults to a no-op implementation.
func NewOrderService(orders OrderRepository, inventory InventoryCoordinator) *OrderService {
	if inventory == nil {
		noop := NewNoopInventoryCoordinator()
		inventory = noop
	}

	return &OrderService{
		orders:    orders,
		inventory: inventory,
	}
}

// CreateOrder handles the business logic of initiating a new purchase.
// It reserves inventory synchronously before writing to the database, ensuring
// that stock cannot be oversold even under high concurrency.
func (s *OrderService) CreateOrder(ctx context.Context, input CreateOrderInput) (CreateOrderResult, error) {
	if s == nil || s.orders == nil || s.inventory == nil {
		return CreateOrderResult{}, fmt.Errorf("order service is not configured")
	}

	normalized, err := normalizeCreateOrderInput(input)
	if err != nil {
		return CreateOrderResult{}, err
	}

	existingOrder, err := s.orders.GetOrderByIdempotencyKey(ctx, normalized.TenantID, normalized.IdempotencyKey)
	switch {
	case err == nil:
		if releaseErr := s.inventory.Release(ctx, inventoryReservationFromOrder(existingOrder)); releaseErr != nil {
			return CreateOrderResult{}, fmt.Errorf("release inventory on idempotent create retry: %w: %v", ErrInventoryUnavailable, releaseErr)
		}
		return CreateOrderResult{Order: existingOrder, Created: false}, nil
	case !errors.Is(err, sql.ErrNoRows):
		return CreateOrderResult{}, fmt.Errorf("lookup order by idempotency key: %w", err)
	}

	reservation := InventoryReservation{
		TenantID:       normalized.TenantID,
		UserID:         normalized.UserID,
		TotalCents:     normalized.TotalCents,
		Currency:       normalized.Currency,
		IdempotencyKey: normalized.IdempotencyKey,
	}

	if err := s.inventory.Reserve(ctx, reservation); err != nil {
		return CreateOrderResult{}, fmt.Errorf("reserve inventory: %w: %v", ErrInventoryUnavailable, err)
	}

	order := models.Order{
		TenantID:       normalized.TenantID,
		UserID:         normalized.UserID,
		Status:         models.OrderStatusPending,
		TotalCents:     normalized.TotalCents,
		Currency:       normalized.Currency,
		IdempotencyKey: normalized.IdempotencyKey,
	}

	if err := s.orders.CreateOrder(ctx, &order); err != nil {
		if errors.Is(err, db.ErrDuplicateValue) {
			existingOrder, lookupErr := s.orders.GetOrderByIdempotencyKey(ctx, normalized.TenantID, normalized.IdempotencyKey)
			if releaseErr := s.inventory.Release(ctx, reservation); releaseErr != nil {
				return CreateOrderResult{}, fmt.Errorf("release inventory after duplicate order: %w: %v", ErrInventoryUnavailable, releaseErr)
			}
			if lookupErr == nil {
				return CreateOrderResult{Order: existingOrder, Created: false}, nil
			}
			return CreateOrderResult{}, fmt.Errorf("load duplicate order by idempotency key: %w", lookupErr)
		}

		if releaseErr := s.inventory.Release(ctx, reservation); releaseErr != nil {
			return CreateOrderResult{}, fmt.Errorf("create order: %w; rollback inventory reservation: %w: %v", err, ErrInventoryUnavailable, releaseErr)
		}
		return CreateOrderResult{}, fmt.Errorf("create order: %w", err)
	}

	return CreateOrderResult{
		Order:   order,
		Created: true,
	}, nil
}

// TransitionOrder safely moves an order from one state to another (e.g. Pending -> Paid).
// It coordinates state-specific side effects, such as releasing inventory holds
// when an order is cancelled or marked as failed.
func (s *OrderService) TransitionOrder(ctx context.Context, req TransitionOrderRequest) (models.Order, error) {
	if s == nil || s.orders == nil || s.inventory == nil {
		return models.Order{}, fmt.Errorf("order service is not configured")
	}

	if req.TenantID <= 0 || req.OrderID <= 0 {
		return models.Order{}, ErrInvalidOrder
	}
	if err := validateTransitionTarget(req.Status); err != nil {
		return models.Order{}, err
	}

	current, err := s.orders.GetOrderByID(ctx, req.TenantID, req.OrderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, ErrOrderNotFound
		}
		return models.Order{}, fmt.Errorf("get order by id: %w", err)
	}

	reservation := InventoryReservation{
		TenantID:       current.TenantID,
		UserID:         current.UserID,
		OrderID:        current.ID,
		TotalCents:     current.TotalCents,
		Currency:       current.Currency,
		IdempotencyKey: current.IdempotencyKey,
	}

	if current.Status == req.Status {
		if shouldRetryReleaseForSameStatus(current.Status) {
			if err := s.inventory.Release(ctx, reservation); err != nil {
				return current, fmt.Errorf("release inventory on idempotent terminal transition: %w: %v", ErrInventoryUnavailable, err)
			}
		}
		return current, nil
	}
	if !canTransitionOrder(current.Status, req.Status) {
		return models.Order{}, ErrInvalidStatusTransition
	}

	if shouldReserveForTransition(current.Status, req.Status) {
		if err := s.inventory.Reserve(ctx, reservation); err != nil {
			return models.Order{}, fmt.Errorf("reserve inventory: %w: %v", ErrInventoryUnavailable, err)
		}
	}

	updated, err := s.orders.UpdateOrderStatus(ctx, req.TenantID, req.OrderID, req.Status)
	if err != nil {
		if shouldReserveForTransition(current.Status, req.Status) {
			if releaseErr := s.inventory.Release(ctx, reservation); releaseErr != nil {
				return models.Order{}, fmt.Errorf("update order status: %w; rollback inventory reservation: %w: %v", err, ErrInventoryUnavailable, releaseErr)
			}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, ErrOrderNotFound
		}
		return models.Order{}, fmt.Errorf("update order status: %w", err)
	}

	// Release happens after persistence so the system does not oversell inventory if the DB update fails.
	if shouldReleaseForTransition(current.Status, req.Status) {
		if err := s.inventory.Release(ctx, reservation); err != nil {
			return updated, fmt.Errorf("release inventory: %w: %v", ErrInventoryUnavailable, err)
		}
	}

	return updated, nil
}

func inventoryReservationFromOrder(order models.Order) InventoryReservation {
	return InventoryReservation{
		TenantID:       order.TenantID,
		UserID:         order.UserID,
		OrderID:        order.ID,
		TotalCents:     order.TotalCents,
		Currency:       order.Currency,
		IdempotencyKey: order.IdempotencyKey,
	}
}
