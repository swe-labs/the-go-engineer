// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package workers

import (
	"context"
	"fmt"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/events"
	paymentflow "github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/payment"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/services"
)

type PaymentWorkflow interface {
	ProcessPayment(ctx context.Context, job paymentflow.Job) (services.ProcessPaymentResult, error)
}

type PaymentProcessor struct {
	Workflow PaymentWorkflow
}

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
