// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package workers

import (
	"context"
	"fmt"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/events"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/services"
)

type OrderWorkflow interface {
	TransitionOrder(ctx context.Context, req services.TransitionOrderRequest) (models.Order, error)
}

type OrderProcessor struct {
	Workflow OrderWorkflow
}

func (p OrderProcessor) Handle(ctx context.Context, event events.Event) error {
	if event.Type != events.TypeOrderStatusChanged {
		return nil
	}
	if p.Workflow == nil {
		return fmt.Errorf("order processor workflow is not configured")
	}
	if event.TenantID <= 0 || event.OrderID <= 0 || event.OrderStatus == "" {
		return events.ErrInvalidEvent
	}

	_, err := p.Workflow.TransitionOrder(ctx, services.TransitionOrderRequest{
		TenantID: event.TenantID,
		OrderID:  event.OrderID,
		Status:   event.OrderStatus,
	})
	if err != nil {
		return fmt.Errorf("process order status event: %w", err)
	}

	return nil
}
