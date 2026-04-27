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
	ErrGatewayTimeout     = errors.New("payment gateway timeout")
	ErrGatewayUnavailable = errors.New("payment gateway unavailable")
)

// Job is the tenant-scoped payment work item used by services and workers.
type Job struct {
	TenantID          int64
	OrderID           int64
	ProviderReference string
	AmountCents       int64
}

type GatewayRequest struct {
	TenantID          int64
	OrderID           int64
	ProviderReference string
	AmountCents       int64
}

type GatewayResult struct {
	Status        models.PaymentStatus
	FailureReason string
}

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

func NewSimulatedGateway() SimulatedGateway {
	return SimulatedGateway{
		Status: models.PaymentStatusSettled,
	}
}

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
