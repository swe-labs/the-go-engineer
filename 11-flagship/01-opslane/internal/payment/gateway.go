// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package payment

import (
	"context"
	"errors"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
)

var (
	// ErrGatewayTimeout is returned when the external payment provider does not respond within the deadline.
	ErrGatewayTimeout = errors.New("payment gateway timeout")
	// ErrGatewayUnavailable is returned when the external provider is experiencing an outage.
	ErrGatewayUnavailable = errors.New("payment gateway unavailable")
)

// Job is the tenant-scoped payment work item used by services and workers.
type Job struct {
	TenantID          int64
	OrderID           int64
	ProviderReference string
	AmountCents       int64
}

// GatewayRequest contains the minimal tenant-scoped information required to process a charge
// against an external provider (e.g., Stripe or Adyen).
type GatewayRequest struct {
	TenantID          int64
	OrderID           int64
	ProviderReference string
	AmountCents       int64
}

// GatewayResult normalizes the response from an external payment provider into
// internal Opslane statuses.
type GatewayResult struct {
	Status        models.PaymentStatus
	FailureReason string
}

// Gateway defines the contract for communicating with external payment processors.
// In a real application, implementations of this interface would wrap the Stripe or Adyen SDKs.
type Gateway interface {
	Charge(ctx context.Context, req GatewayRequest) (GatewayResult, error)
}

// SimulatedGateway gives the flagship a deterministic gateway boundary before real providers exist.
type SimulatedGateway struct {
	Delay         time.Duration
	Status        models.PaymentStatus
	FailureReason string
	Err           error
}

// NewSimulatedGateway creates a test-friendly Gateway that always succeeds immediately.
func NewSimulatedGateway() SimulatedGateway {
	return SimulatedGateway{
		Status: models.PaymentStatusSettled,
	}
}

// Charge simulates an external network call, applying the configured delay and returning the mock status.
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
