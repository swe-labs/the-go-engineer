package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/db"
	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/models"
	paymentflow "github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/payment"
)

func TestPaymentServiceSettlesPaymentAndMarksOrderPaid(t *testing.T) {
	t.Parallel()

	orders := newStubOrderRepository()
	orders.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusPending,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	payments := newStubPaymentRepository()
	service := NewPaymentService(payments, orders, NewOrderService(orders, &stubInventoryCoordinator{}), &stubPaymentGateway{
		results: []paymentflow.GatewayResult{{Status: models.PaymentStatusSettled}},
	}, PaymentServiceOptions{GatewayTimeout: time.Second, MaxAttempts: 1})

	result, err := service.ProcessPayment(context.Background(), paymentflow.Job{
		TenantID:          7,
		OrderID:           101,
		ProviderReference: "pay_123",
		AmountCents:       2500,
	})
	if err != nil {
		t.Fatalf("ProcessPayment returned error: %v", err)
	}

	if !result.Created {
		t.Fatal("Created = false, want true")
	}
	if result.Payment.Status != models.PaymentStatusSettled {
		t.Fatalf("payment status = %q, want settled", result.Payment.Status)
	}
	order, err := orders.GetOrderByID(context.Background(), 7, 101)
	if err != nil {
		t.Fatalf("GetOrderByID returned error: %v", err)
	}
	if order.Status != models.OrderStatusPaid {
		t.Fatalf("order status = %q, want paid", order.Status)
	}
}

func TestPaymentServiceReturnsExistingPaymentForDuplicateProviderReference(t *testing.T) {
	t.Parallel()

	orders := newStubOrderRepository()
	orders.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusPaid,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	payments := newStubPaymentRepository()
	payments.seedPayment(models.Payment{
		ID:                501,
		TenantID:          7,
		OrderID:           101,
		Status:            models.PaymentStatusSettled,
		ProviderReference: "pay_123",
		AmountCents:       2500,
	})
	gateway := &stubPaymentGateway{}
	service := NewPaymentService(payments, orders, NewOrderService(orders, &stubInventoryCoordinator{}), gateway, PaymentServiceOptions{})

	result, err := service.ProcessPayment(context.Background(), paymentflow.Job{
		TenantID:          7,
		OrderID:           101,
		ProviderReference: "pay_123",
		AmountCents:       2500,
	})
	if err != nil {
		t.Fatalf("ProcessPayment returned error: %v", err)
	}

	if result.Created {
		t.Fatal("Created = true, want false")
	}
	if result.Payment.ID != 501 {
		t.Fatalf("payment id = %d, want 501", result.Payment.ID)
	}
	if gateway.calls != 0 {
		t.Fatalf("gateway calls = %d, want 0", gateway.calls)
	}
}

func TestPaymentServiceKeepsPendingPaymentWhenGatewayTimesOut(t *testing.T) {
	t.Parallel()

	orders := newStubOrderRepository()
	orders.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusPending,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	payments := newStubPaymentRepository()
	service := NewPaymentService(payments, orders, NewOrderService(orders, &stubInventoryCoordinator{}), &stubPaymentGateway{
		err: paymentflow.ErrGatewayTimeout,
	}, PaymentServiceOptions{GatewayTimeout: time.Second, MaxAttempts: 1})

	result, err := service.ProcessPayment(context.Background(), paymentflow.Job{
		TenantID:          7,
		OrderID:           101,
		ProviderReference: "pay_timeout",
		AmountCents:       2500,
	})
	if !errors.Is(err, paymentflow.ErrGatewayTimeout) {
		t.Fatalf("ProcessPayment error = %v, want ErrGatewayTimeout", err)
	}
	if result.Payment.Status != models.PaymentStatusPending {
		t.Fatalf("payment status = %q, want pending", result.Payment.Status)
	}

	order, err := orders.GetOrderByID(context.Background(), 7, 101)
	if err != nil {
		t.Fatalf("GetOrderByID returned error: %v", err)
	}
	if order.Status != models.OrderStatusProcessing {
		t.Fatalf("order status = %q, want processing", order.Status)
	}
}

func TestPaymentServiceReconcilePaymentSettlesTimedOutPayment(t *testing.T) {
	t.Parallel()

	orders := newStubOrderRepository()
	orders.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusProcessing,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	payments := newStubPaymentRepository()
	payments.seedPayment(models.Payment{
		ID:                501,
		TenantID:          7,
		OrderID:           101,
		Status:            models.PaymentStatusPending,
		ProviderReference: "pay_timeout",
		AmountCents:       2500,
	})
	service := NewPaymentService(payments, orders, NewOrderService(orders, &stubInventoryCoordinator{}), &stubPaymentGateway{}, PaymentServiceOptions{})

	updated, err := service.ReconcilePayment(context.Background(), ReconcilePaymentInput{
		TenantID:          7,
		ProviderReference: "pay_timeout",
		Status:            models.PaymentStatusSettled,
	})
	if err != nil {
		t.Fatalf("ReconcilePayment returned error: %v", err)
	}
	if updated.Status != models.PaymentStatusSettled {
		t.Fatalf("payment status = %q, want settled", updated.Status)
	}

	order, err := orders.GetOrderByID(context.Background(), 7, 101)
	if err != nil {
		t.Fatalf("GetOrderByID returned error: %v", err)
	}
	if order.Status != models.OrderStatusPaid {
		t.Fatalf("order status = %q, want paid", order.Status)
	}
}

func TestPaymentServiceRejectsProviderReferenceReusedForDifferentOrder(t *testing.T) {
	t.Parallel()

	orders := newStubOrderRepository()
	orders.seedOrder(models.Order{
		ID:             101,
		TenantID:       7,
		UserID:         42,
		Status:         models.OrderStatusPending,
		TotalCents:     2500,
		Currency:       "USD",
		IdempotencyKey: "checkout-123",
	})
	payments := newStubPaymentRepository()
	payments.seedPayment(models.Payment{
		ID:                501,
		TenantID:          7,
		OrderID:           202,
		Status:            models.PaymentStatusSettled,
		ProviderReference: "pay_123",
		AmountCents:       2500,
	})
	service := NewPaymentService(payments, orders, NewOrderService(orders, &stubInventoryCoordinator{}), &stubPaymentGateway{}, PaymentServiceOptions{})

	_, err := service.ProcessPayment(context.Background(), paymentflow.Job{
		TenantID:          7,
		OrderID:           101,
		ProviderReference: "pay_123",
		AmountCents:       2500,
	})
	if !errors.Is(err, ErrInvalidPayment) {
		t.Fatalf("ProcessPayment error = %v, want ErrInvalidPayment", err)
	}
}

type stubPaymentGateway struct {
	results []paymentflow.GatewayResult
	err     error
	calls   int
}

func (g *stubPaymentGateway) Charge(context.Context, paymentflow.GatewayRequest) (paymentflow.GatewayResult, error) {
	g.calls++
	if g.err != nil {
		return paymentflow.GatewayResult{}, g.err
	}
	if len(g.results) == 0 {
		return paymentflow.GatewayResult{Status: models.PaymentStatusSettled}, nil
	}

	result := g.results[0]
	g.results = g.results[1:]
	return result, nil
}

type stubPaymentRepository struct {
	payments    map[string]models.Payment
	nextID      int64
	createErr   error
	updateErr   error
	createCalls int
	updateCalls int
}

func newStubPaymentRepository() *stubPaymentRepository {
	return &stubPaymentRepository{
		payments: make(map[string]models.Payment),
		nextID:   501,
	}
}

func (r *stubPaymentRepository) CreatePayment(_ context.Context, payment *models.Payment) error {
	r.createCalls++
	if r.createErr != nil {
		return r.createErr
	}

	key := r.makeKey(payment.TenantID, payment.ProviderReference)
	if _, ok := r.payments[key]; ok {
		return db.ErrDuplicateValue
	}

	payment.ID = r.nextID
	r.nextID++
	now := time.Now().UTC()
	payment.CreatedAt = now
	payment.UpdatedAt = now
	r.seedPayment(*payment)
	return nil
}

func (r *stubPaymentRepository) GetPaymentByProviderReference(_ context.Context, tenantID int64, providerReference string) (models.Payment, error) {
	payment, ok := r.payments[r.makeKey(tenantID, providerReference)]
	if !ok {
		return models.Payment{}, sql.ErrNoRows
	}

	return payment, nil
}

func (r *stubPaymentRepository) UpdatePaymentStatus(_ context.Context, tenantID int64, providerReference string, status models.PaymentStatus, failureReason string) (models.Payment, error) {
	r.updateCalls++
	if r.updateErr != nil {
		return models.Payment{}, r.updateErr
	}

	key := r.makeKey(tenantID, providerReference)
	payment, ok := r.payments[key]
	if !ok {
		return models.Payment{}, sql.ErrNoRows
	}

	payment.Status = status
	payment.FailureReason = failureReason
	payment.UpdatedAt = time.Now().UTC()
	r.seedPayment(payment)
	return payment, nil
}

func (r *stubPaymentRepository) seedPayment(payment models.Payment) {
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = time.Now().UTC()
	}
	if payment.UpdatedAt.IsZero() {
		payment.UpdatedAt = payment.CreatedAt
	}
	r.payments[r.makeKey(payment.TenantID, payment.ProviderReference)] = payment
}

func (r *stubPaymentRepository) makeKey(tenantID int64, providerReference string) string {
	return fmt.Sprintf("%d::%s", tenantID, providerReference)
}
