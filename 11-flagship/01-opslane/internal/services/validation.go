// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package services

import (
	"errors"
	"strings"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

var (
	ErrInvalidOrder            = errors.New("invalid order")
	ErrOrderNotFound           = errors.New("order not found")
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrInventoryUnavailable    = errors.New("inventory unavailable")
)

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

func validateTransitionTarget(status models.OrderStatus) error {
	if !isKnownOrderStatus(status) {
		return ErrInvalidStatusTransition
	}

	return nil
}

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

func shouldReserveForTransition(current, next models.OrderStatus) bool {
	return current == models.OrderStatusFailed && next == models.OrderStatusProcessing
}

func shouldReleaseForTransition(current, next models.OrderStatus) bool {
	switch current {
	case models.OrderStatusPending, models.OrderStatusProcessing:
		return next == models.OrderStatusFailed || next == models.OrderStatusCancelled
	default:
		return false
	}
}

func shouldRetryReleaseForSameStatus(status models.OrderStatus) bool {
	return status == models.OrderStatusFailed || status == models.OrderStatusCancelled
}
