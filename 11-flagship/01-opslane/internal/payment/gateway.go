// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

// Package payment provides the payment gateway abstraction for the Opslane backend.
// It defines the Gateway interface for external payment providers and includes
// a simulated gateway for testing and development.
package payment

import (
	"context"
	"errors"
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/models"
)

var (
	// ErrGatewayTimeout (Error): returned when the external payment provider does not respond within the deadline
	ErrGatewayTimeout = errors.New("payment gateway timeout")
	// ErrGatewayUnavailable (Error): returned when the external provider is experiencing an outage
	ErrGatewayUnavailable = errors.New("payment gateway unavailable")
)

// Job (Struct): tenant-scoped payment work item used by services and workers
type Job struct {
	TenantID          int64
	OrderID           int64
	ProviderReference string
	AmountCents       int64
}

// GatewayRequest (Struct): tenant-scoped information required to process a charge against an external provider
type GatewayRequest struct {
	TenantID          int64
	OrderID           int64
	ProviderReference string
	AmountCents       int64
}

// GatewayResult (Struct): normalizes the response from an external payment provider into internal statuses
type GatewayResult struct {
	Status        models.PaymentStatus
	FailureReason string
}

// Gateway (Interface): contract for communicating with external payment processors
type Gateway interface {
	Charge(ctx context.Context, req GatewayRequest) (GatewayResult, error)
}

// SimulatedGateway (Struct): deterministic gateway boundary for development and testing before real providers exist
type SimulatedGateway struct {
	Delay         time.Duration
	Status        models.PaymentStatus
	FailureReason string
	Err           error
}

// NewSimulatedGateway (Constructor): creates a test-friendly Gateway that always succeeds immediately
func NewSimulatedGateway() SimulatedGateway {
	return SimulatedGateway{
		Status: models.PaymentStatusSettled,
	}
}

// Charge (Method): simulates an external network call with configured delay and mock status/error
func (g SimulatedGateway) Charge(ctx context.Context, _ GatewayRequest) (GatewayResult, error) {
	if g.Delay > 0 {
		timer := time.NewTimer(g.Delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return GatewayResult{}, ErrGatewayTimeout
		case <-timer.C:
		}
	}

	if g.Err != nil {
		return GatewayResult{}, g.Err
	}

	status := g.Status
	if status == "" {
		status = models.PaymentStatusSettled
	}

	return GatewayResult{
		Status:        status,
		FailureReason: g.FailureReason,
	}, nil
}
