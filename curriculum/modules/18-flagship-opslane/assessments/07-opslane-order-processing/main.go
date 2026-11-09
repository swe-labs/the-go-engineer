package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var (
	ErrInvalidStateTransition = errors.New("invalid order state transition")
	ErrOrderNotFound          = errors.New("order not found")
	ErrPaymentFailed          = errors.New("payment execution failed")
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusCompleted  OrderStatus = "COMPLETED"
	StatusFailed     OrderStatus = "FAILED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

// AllowedTransitions defines valid state progression rules.
var AllowedTransitions = map[OrderStatus][]OrderStatus{
	StatusPending:    {StatusProcessing, StatusCancelled},
	StatusProcessing: {StatusCompleted, StatusFailed},
	StatusFailed:     {StatusProcessing, StatusCancelled}, // Retry capability
	StatusCompleted:  {},                                  // Terminal
	StatusCancelled:  {},                                  // Terminal
}

func CanTransition(from, to OrderStatus) bool {
	allowed, ok := AllowedTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// PaymentGateway interface
type PaymentGateway interface {
	ProcessPayment(ctx context.Context, orderID string, amount float64) error
}

// OrderProcessor coordinates transactional order workflows.
type OrderProcessor struct {
	db      *sql.DB
	gateway PaymentGateway
}

func NewOrderProcessor(db *sql.DB, gateway PaymentGateway) (*OrderProcessor, error) {
	schema := `
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		tenant_id TEXT NOT NULL,
		amount REAL NOT NULL,
		status TEXT NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id TEXT NOT NULL,
		action TEXT NOT NULL,
		status TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}
	return &OrderProcessor{db: db, gateway: gateway}, nil
}

// ProcessOrder atomically advances an order from PENDING to COMPLETED, or rolls back on payment failure.
func (op *OrderProcessor) ProcessOrder(ctx context.Context, orderID string) error {
	tx, err := op.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback() // Safe rollback if not committed

	// 1. Fetch current order status
	var currentStatus OrderStatus
	var amount float64
	err = tx.QueryRowContext(ctx, "SELECT status, amount FROM orders WHERE id = ?", orderID).Scan(&currentStatus, &amount)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrOrderNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to query order: %w", err)
	}

	// 2. Validate state machine invariant
	if !CanTransition(currentStatus, StatusProcessing) {
		return fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidStateTransition, currentStatus, StatusProcessing)
	}

	now := time.Now().UTC()

	// 3. Mark as PROCESSING in transaction
	_, err = tx.ExecContext(ctx, "UPDATE orders SET status = ?, updated_at = ? WHERE id = ?", StatusProcessing, now, orderID)
	if err != nil {
		return fmt.Errorf("failed to mark processing: %w", err)
	}

	// 4. Execute external payment gateway
	if err := op.gateway.ProcessPayment(ctx, orderID, amount); err != nil {
		// If payment fails, explicitly roll back transaction so order remains in PENDING status
		return fmt.Errorf("%w: %v", ErrPaymentFailed, err)
	}

	// 5. If payment succeeds, transition to COMPLETED
	_, err = tx.ExecContext(ctx, "UPDATE orders SET status = ?, updated_at = ? WHERE id = ?", StatusCompleted, now, orderID)
	if err != nil {
		return fmt.Errorf("failed to mark completed: %w", err)
	}

	// 6. Write audit log entry in the same atomic transaction
	auditQuery := `INSERT INTO audit_logs (order_id, action, status, timestamp) VALUES (?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, auditQuery, orderID, "ORDER_PAYMENT_PROCESSED", StatusCompleted, now)
	if err != nil {
		return fmt.Errorf("failed to record audit log: %w", err)
	}

	// 7. Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit order transaction: %w", err)
	}

	return nil
}

type MockPaymentGateway struct {
	ShouldFail bool
}

func (m *MockPaymentGateway) ProcessPayment(ctx context.Context, orderID string, amount float64) error {
	if m.ShouldFail {
		return errors.New("gateway declined: insufficient funds")
	}
	return nil
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()

	processor, err := NewOrderProcessor(db, &MockPaymentGateway{ShouldFail: false})
	if err != nil {
		log.Fatalf("Init error: %v", err)
	}

	_, _ = db.Exec("INSERT INTO orders (id, tenant_id, amount, status, updated_at) VALUES ('ord-demo', 't-1', 100.0, 'PENDING', ?)", time.Now())
	if err := processor.ProcessOrder(context.Background(), "ord-demo"); err != nil {
		log.Fatalf("Process error: %v", err)
	}

	fmt.Println("Opslane Order Processing transaction completed successfully.")
}
