package workers

import (
	"context"
	"errors"
	"testing"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/events"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
	paymentflow "github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/payment"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/services"
)

func TestOrderProcessorTransitionsOrderStatus(t *testing.T) {
	t.Parallel()

	workflow := &stubOrderWorkflow{}
	processor := OrderProcessor{Workflow: workflow}

	err := processor.Handle(context.Background(), events.Event{
		Type:        events.TypeOrderStatusChanged,
		TenantID:    7,
		OrderID:     101,
		OrderStatus: models.OrderStatusPaid,
	})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if workflow.last.TenantID != 7 || workflow.last.OrderID != 101 || workflow.last.Status != models.OrderStatusPaid {
		t.Fatalf("transition = %+v, want tenant 7 order 101 paid", workflow.last)
	}
}

func TestPaymentProcessorStartsPaymentWorkflow(t *testing.T) {
	t.Parallel()

	workflow := &stubPaymentWorkflow{}
	processor := PaymentProcessor{Workflow: workflow}

	err := processor.Handle(context.Background(), events.Event{
		Type:              events.TypePaymentRequested,
		TenantID:          7,
		OrderID:           101,
		ProviderReference: "pay_123",
		AmountCents:       2500,
	})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if workflow.last.TenantID != 7 || workflow.last.OrderID != 101 || workflow.last.ProviderReference != "pay_123" {
		t.Fatalf("payment job = %+v, want tenant 7 order 101 pay_123", workflow.last)
	}
}

func TestNotificationWorkerSendsNotification(t *testing.T) {
	t.Parallel()

	sink := &stubNotificationSink{}
	worker := NotificationWorker{Sink: sink}

	err := worker.Handle(context.Background(), events.Event{
		ID:       "evt_1",
		Type:     events.TypeNotificationRequested,
		TenantID: 7,
		Metadata: map[string]string{
			"channel": "email",
			"subject": "Payment settled",
			"body":    "Your payment succeeded.",
		},
	})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if sink.last.TenantID != 7 || sink.last.Channel != "email" || sink.last.Subject != "Payment settled" {
		t.Fatalf("notification = %+v, want tenant 7 email subject", sink.last)
	}
}

func TestNotificationWorkerSurfacesSinkError(t *testing.T) {
	t.Parallel()

	sinkErr := errors.New("smtp down")
	worker := NotificationWorker{Sink: &stubNotificationSink{err: sinkErr}}

	err := worker.Handle(context.Background(), events.Event{
		ID:       "evt_1",
		Type:     events.TypeNotificationRequested,
		TenantID: 7,
		Metadata: map[string]string{
			"channel": "email",
			"subject": "Payment settled",
		},
	})
	if !errors.Is(err, sinkErr) {
		t.Fatalf("Handle error = %v, want sinkErr", err)
	}
}

type stubOrderWorkflow struct {
	last services.TransitionOrderRequest
	err  error
}

func (s *stubOrderWorkflow) TransitionOrder(_ context.Context, req services.TransitionOrderRequest) (models.Order, error) {
	s.last = req
	return models.Order{ID: req.OrderID, TenantID: req.TenantID, Status: req.Status}, s.err
}

type stubPaymentWorkflow struct {
	last paymentflow.Job
	err  error
}

func (s *stubPaymentWorkflow) ProcessPayment(_ context.Context, job paymentflow.Job) (services.ProcessPaymentResult, error) {
	s.last = job
	return services.ProcessPaymentResult{}, s.err
}

type stubNotificationSink struct {
	last Notification
	err  error
}

func (s *stubNotificationSink) Send(_ context.Context, notification Notification) error {
	s.last = notification
	return s.err
}
