package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

func setupTestProcessor(t *testing.T, gateway PaymentGateway) (*OrderProcessor, *sql.DB) {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	processor, err := NewOrderProcessor(db, gateway)
	if err != nil {
		t.Fatalf("failed to init processor: %v", err)
	}

	return processor, db
}

func TestOrderProcessing_Success(t *testing.T) {
	gateway := &MockPaymentGateway{ShouldFail: false}
	processor, db := setupTestProcessor(t, gateway)
	defer db.Close()

	// Seed order in PENDING status
	now := time.Now().UTC()
	_, err := db.Exec("INSERT INTO orders (id, tenant_id, amount, status, updated_at) VALUES ('ord-success-1', 'tenant-acme', 250.0, 'PENDING', ?)", now)
	if err != nil {
		t.Fatalf("failed to seed order: %v", err)
	}

	err = processor.ProcessOrder(context.Background(), "ord-success-1")
	if err != nil {
		t.Fatalf("expected order processing to succeed, got: %v", err)
	}

	// Verify order status is COMPLETED
	var status OrderStatus
	err = db.QueryRow("SELECT status FROM orders WHERE id = 'ord-success-1'").Scan(&status)
	if err != nil || status != StatusCompleted {
		t.Errorf("expected order status COMPLETED, got %s (err: %v)", status, err)
	}

	// Verify audit log entry was written
	var count int
	err = db.QueryRow("SELECT COUNT(1) FROM audit_logs WHERE order_id = 'ord-success-1' AND action = 'ORDER_PAYMENT_PROCESSED'").Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("expected 1 audit log entry, got %d", count)
	}
}

func TestOrderProcessing_PaymentFailureRollback(t *testing.T) {
	// Core practice task requirement:
	// Trigger failure during payment processing, verify transaction rolls back,
	// and confirm that the order remains in PENDING status.
	gateway := &MockPaymentGateway{ShouldFail: true}
	processor, db := setupTestProcessor(t, gateway)
	defer db.Close()

	// Seed order in PENDING status
	now := time.Now().UTC()
	_, err := db.Exec("INSERT INTO orders (id, tenant_id, amount, status, updated_at) VALUES ('ord-fail-2', 'tenant-acme', 500.0, 'PENDING', ?)", now)
	if err != nil {
		t.Fatalf("failed to seed order: %v", err)
	}

	err = processor.ProcessOrder(context.Background(), "ord-fail-2")
	if !errors.Is(err, ErrPaymentFailed) {
		t.Fatalf("expected ErrPaymentFailed, got: %v", err)
	}

	// Verify order rolled back to original PENDING status
	var status OrderStatus
	err = db.QueryRow("SELECT status FROM orders WHERE id = 'ord-fail-2'").Scan(&status)
	if err != nil {
		t.Fatalf("failed to query order status: %v", err)
	}
	if status != StatusPending {
		t.Fatalf("TRANSACTION ROLLBACK FAILURE: expected order to remain PENDING, but got %s", status)
	}

	// Verify NO audit log was committed
	var count int
	_ = db.QueryRow("SELECT COUNT(1) FROM audit_logs WHERE order_id = 'ord-fail-2'").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 audit logs due to rollback, got %d", count)
	}
}

func TestOrderProcessing_InvalidTransition(t *testing.T) {
	gateway := &MockPaymentGateway{ShouldFail: false}
	processor, db := setupTestProcessor(t, gateway)
	defer db.Close()

	// Seed order already in COMPLETED status (terminal)
	_, _ = db.Exec("INSERT INTO orders (id, tenant_id, amount, status, updated_at) VALUES ('ord-done-3', 'tenant-acme', 100.0, 'COMPLETED', ?)", time.Now())

	err := processor.ProcessOrder(context.Background(), "ord-done-3")
	if !errors.Is(err, ErrInvalidStateTransition) {
		t.Fatalf("expected ErrInvalidStateTransition for completed order, got: %v", err)
	}
}
