package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type DedupStore struct {
	mu     sync.Mutex
	seen   map[string]bool
	onMiss func(string)
	onDup  func(string)
}

func NewDedupStore(onMiss, onDup func(string)) *DedupStore {
	return &DedupStore{
		seen:   make(map[string]bool),
		onMiss: onMiss,
		onDup:  onDup,
	}
}

func (d *DedupStore) TryProcess(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.seen[key] {
		if d.onDup != nil {
			d.onDup(key)
		}
		return false
	}
	d.seen[key] = true
	if d.onMiss != nil {
		d.onMiss(key)
	}
	return true
}

func (d *DedupStore) Seen(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.seen[key]
}

func (d *DedupStore) Count() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Jitter      bool
}

type Result struct {
	Attempts int
	Success  bool
	Err      error
}

func RetryWithBackoff(cfg RetryConfig, op func() error) Result {
	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err := op()
		if err == nil {
			return Result{Attempts: attempt, Success: true}
		}
		lastErr = err
		if attempt == cfg.MaxAttempts {
			break
		}
		delay := cfg.BaseDelay * (1 << uint(attempt-1))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
		if cfg.Jitter {
			jitter := time.Duration(rand.Int63n(int64(delay))) / 2
			delay = delay/2 + jitter
		}
		time.Sleep(delay)
	}
	return Result{Attempts: cfg.MaxAttempts, Success: false, Err: lastErr}
}

var failCounter int
var failMu sync.Mutex

func flakyOperation(succeedAfter int) func() error {
	return func() error {
		failMu.Lock()
		defer failMu.Unlock()
		failCounter++
		if failCounter < succeedAfter {
			return fmt.Errorf("transient error (attempt %d)", failCounter)
		}
		return nil
	}
}

func main() {
	idempotencyKeys := []string{"payment-1", "refund-2", "payment-1"}
	processed := 0
	deduped := 0
	store := NewDedupStore(
		func(k string) { processed++ },
		func(k string) { deduped++ },
	)

	for _, key := range idempotencyKeys {
		store.TryProcess(key)
	}
	fmt.Printf("Processed: %d, Deduped: %d\n", processed, deduped)

	failCounter = 0
	cfg := RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    100 * time.Millisecond,
		Jitter:      false,
	}
	result := RetryWithBackoff(cfg, flakyOperation(3))
	fmt.Printf("Retry result: success=%v attempts=%d err=%v\n", result.Success, result.Attempts, result.Err)
}
