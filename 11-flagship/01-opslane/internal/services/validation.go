// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package services

import (
	"errors"
	"strings"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/models"
)

var (
	// ErrInvalidOrder (Error): returned when order input validation fails
	ErrInvalidOrder = errors.New("invalid order")
	// ErrOrderNotFound (Error): returned when a requested order does not exist
	ErrOrderNotFound = errors.New("order not found")
	// ErrInvalidStatusTransition (Error): returned when an order status change is not allowed
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	// ErrInventoryUnavailable (Error): returned when inventory reservation fails
	ErrInventoryUnavailable = errors.New("inventory unavailable")
)

// normalizeCreateOrderInput (Function): validates and normalizes order creation input
func normalizeCreateOrderInput(input CreateOrderInput) (CreateOrderInput, error) {
	input.Currency = strings.ToUpper(strings.TrimSpace(input.Currency))
	input.IdempotencyKey = strings.TrimSpace(input.IdempotencyKey)

	if input.TenantID <= 0 || input.UserID <= 0 || input.TotalCents <= 0 || input.IdempotencyKey == "" {
		return CreateOrderInput{}, ErrInvalidOrder
	}

	if len(input.Currency) != 3 {
		return CreateOrderInput{}, ErrInvalidOrder
	}

	return input, nil
}

// validateTransitionTarget (Function): validates that a target order status is a known value
func validateTransitionTarget(status models.OrderStatus) error {
	if !isKnownOrderStatus(status) {
		return ErrInvalidStatusTransition
	}

	return nil
}

// isKnownOrderStatus (Function): checks if an order status is a recognized enum value
func isKnownOrderStatus(status models.OrderStatus) bool {
	switch status {
	case models.OrderStatusPending,
		models.OrderStatusProcessing,
		models.OrderStatusPaid,
		models.OrderStatusFailed,
		models.OrderStatusCancelled:
		return true
	default:
		return false
	}
}

// canTransitionOrder (Function): determines if an order status transition is valid
func canTransitionOrder(current, next models.OrderStatus) bool {
	switch current {
	case models.OrderStatusPending:
		return next == models.OrderStatusProcessing || next == models.OrderStatusFailed || next == models.OrderStatusCancelled
	case models.OrderStatusProcessing:
		return next == models.OrderStatusPaid || next == models.OrderStatusFailed || next == models.OrderStatusCancelled
	case models.OrderStatusFailed:
		return next == models.OrderStatusProcessing || next == models.OrderStatusCancelled
	default:
		return false
	}
}

// shouldReserveForTransition (Function): determines if inventory reservation is needed before a transition
func shouldReserveForTransition(current, next models.OrderStatus) bool {
	return current == models.OrderStatusFailed && next == models.OrderStatusProcessing
}

// shouldReleaseForTransition (Function): determines if inventory release is needed after a transition
func shouldReleaseForTransition(current, next models.OrderStatus) bool {
	switch current {
	case models.OrderStatusPending, models.OrderStatusProcessing:
		return next == models.OrderStatusFailed || next == models.OrderStatusCancelled
	default:
		return false
	}
}

// shouldRetryReleaseForSameStatus (Function): determines if inventory should be released on idempotent same-status retry
func shouldRetryReleaseForSameStatus(status models.OrderStatus) bool {
	return status == models.OrderStatusFailed || status == models.OrderStatusCancelled
}
