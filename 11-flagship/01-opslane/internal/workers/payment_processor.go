// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package workers

import (
	"context"
	"fmt"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/events"
	paymentflow "github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/payment"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/services"
)

// PaymentWorkflow (Interface): processes payment jobs through the payment service
type PaymentWorkflow interface {
	ProcessPayment(ctx context.Context, job paymentflow.Job) (services.ProcessPaymentResult, error)
}

// PaymentProcessor (Struct): handles TypePaymentRequested events by delegating to the payment workflow
type PaymentProcessor struct {
	Workflow PaymentWorkflow
}

// Handle (Method): processes a TypePaymentRequested event through the payment workflow
func (p PaymentProcessor) Handle(ctx context.Context, event events.Event) error {
	if event.Type != events.TypePaymentRequested {
		return nil
	}
	if p.Workflow == nil {
		return fmt.Errorf("payment processor workflow is not configured")
	}

	_, err := p.Workflow.ProcessPayment(ctx, paymentflow.Job{
		TenantID:          event.TenantID,
		OrderID:           event.OrderID,
		ProviderReference: event.ProviderReference,
		AmountCents:       event.AmountCents,
	})
	if err != nil {
		return fmt.Errorf("process payment event: %w", err)
	}

	return nil
}
