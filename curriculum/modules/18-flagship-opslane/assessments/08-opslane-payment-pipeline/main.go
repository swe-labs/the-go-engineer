package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrDuplicateInFlight = errors.New("concurrent payment with same idempotency key in progress")
	ErrCardDeclined      = errors.New("card declined: insufficient funds")
	ErrTransientNetwork  = errors.New("gateway timeout / transient network error")
)

type PaymentReceipt struct {
	TransactionID  string    `json:"transaction_id"`
	OrderID        string    `json:"order_id"`
	IdempotencyKey string    `json:"idempotency_key"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"` // "COMPLETED"
	ProcessedAt    time.Time `json:"processed_at"`
}

type IdempotencyRecord struct {
	Status  string // "IN_PROGRESS", "COMPLETED"
	Receipt *PaymentReceipt
}

// PaymentGatewayAdapter defines external processor interface.
type PaymentGatewayAdapter interface {
	ExecuteCharge(ctx context.Context, orderID string, amount float64) (string, error)
}

// PaymentPipeline manages idempotency and retry mechanics.
type PaymentPipeline struct {
	gateway    PaymentGatewayAdapter
	mu         sync.Mutex
	records    map[string]*IdempotencyRecord
	maxRetries int
	baseDelay  time.Duration
}

func NewPaymentPipeline(gateway PaymentGatewayAdapter, maxRetries int, baseDelay time.Duration) *PaymentPipeline {
	return &PaymentPipeline{
		gateway:    gateway,
		records:    make(map[string]*IdempotencyRecord),
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
	}
}

// Charge executes payment with strict idempotency and backoff retries.
func (pp *PaymentPipeline) Charge(ctx context.Context, idempotencyKey, orderID string, amount float64) (*PaymentReceipt, error) {
	pp.mu.Lock()
	record, exists := pp.records[idempotencyKey]
	if exists {
		if record.Status == "COMPLETED" {
			receipt := record.Receipt
			pp.mu.Unlock()
			// Idempotent cache hit: return original receipt without re-charging
			return receipt, nil
		}
		if record.Status == "IN_PROGRESS" {
			pp.mu.Unlock()
			return nil, ErrDuplicateInFlight
		}
	}

	// Reserve key
	pp.records[idempotencyKey] = &IdempotencyRecord{Status: "IN_PROGRESS"}
	pp.mu.Unlock()

	// Execute with exponential backoff retries for transient errors
	var txID string
	var err error
	delay := pp.baseDelay

	for attempt := 0; attempt <= pp.maxRetries; attempt++ {
		txID, err = pp.gateway.ExecuteCharge(ctx, orderID, amount)
		if err == nil {
			break
		}

		// Fail immediately if error is non-transient
		if !errors.Is(err, ErrTransientNetwork) {
			pp.mu.Lock()
			delete(pp.records, idempotencyKey)
			pp.mu.Unlock()
			return nil, err
		}

		// If transient and retries remain, back off
		if attempt < pp.maxRetries {
			select {
			case <-ctx.Done():
				pp.mu.Lock()
				delete(pp.records, idempotencyKey)
				pp.mu.Unlock()
				return nil, ctx.Err()
			case <-time.After(delay):
				delay *= 2 // Exponential backoff
			}
		}
	}

	if err != nil {
		pp.mu.Lock()
		delete(pp.records, idempotencyKey)
		pp.mu.Unlock()
		return nil, fmt.Errorf("payment failed after retries: %w", err)
	}

	receipt := &PaymentReceipt{
		TransactionID:  txID,
		OrderID:        orderID,
		IdempotencyKey: idempotencyKey,
		Amount:         amount,
		Status:         "COMPLETED",
		ProcessedAt:    time.Now().UTC(),
	}

	pp.mu.Lock()
	pp.records[idempotencyKey] = &IdempotencyRecord{
		Status:  "COMPLETED",
		Receipt: receipt,
	}
	pp.mu.Unlock()

	return receipt, nil
}

// FlakyGatewaySimulator simulates transient errors and call tracking.
type FlakyGatewaySimulator struct {
	CallCount       int32
	TransientErrors int32 // Number of times to fail before succeeding
	PermanentFail   bool
}

func (f *FlakyGatewaySimulator) ExecuteCharge(ctx context.Context, orderID string, amount float64) (string, error) {
	currentCall := atomic.AddInt32(&f.CallCount, 1)

	if f.PermanentFail {
		return "", ErrCardDeclined
	}

	if currentCall <= atomic.LoadInt32(&f.TransientErrors) {
		return "", ErrTransientNetwork
	}

	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "tx-" + hex.EncodeToString(b), nil
}

func main() {
	gateway := &FlakyGatewaySimulator{TransientErrors: 1}
	pipeline := NewPaymentPipeline(gateway, 3, 10*time.Millisecond)

	receipt1, err := pipeline.Charge(context.Background(), "idem-key-100", "ord-99", 199.99)
	if err != nil {
		log.Fatalf("Charge 1 failed: %v", err)
	}
	fmt.Printf("Charge 1 completed: Tx %s\n", receipt1.TransactionID)

	// Second identical charge -> returns same receipt without calling gateway again
	receipt2, err := pipeline.Charge(context.Background(), "idem-key-100", "ord-99", 199.99)
	if err != nil {
		log.Fatalf("Charge 2 failed: %v", err)
	}
	fmt.Printf("Charge 2 (idempotent replay): Tx %s (matches: %t)\n", receipt2.TransactionID, receipt1.TransactionID == receipt2.TransactionID)
}
