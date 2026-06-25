package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type State int

const (
	StateClosed   State = 0
	StateOpen     State = 1
	StateHalfOpen State = 2
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	}
	return "unknown"
}

type CircuitBreaker struct {
	mu              sync.Mutex
	state           State
	failureCount    int
	successCount    int
	threshold       int
	halfOpenMax     int
	cooldown        time.Duration
	lastFailureTime time.Time
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:       StateClosed,
		threshold:   threshold,
		halfOpenMax: 3,
		cooldown:    cooldown,
	}
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

func (cb *CircuitBreaker) Call(ctx context.Context, fn func(context.Context) error) error {
	cb.mu.Lock()

	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailureTime) < cb.cooldown {
			cb.mu.Unlock()
			return errors.New("circuit breaker open")
		}
		cb.state = StateHalfOpen
		cb.successCount = 0
	case StateHalfOpen:
		if cb.successCount >= cb.halfOpenMax {
			cb.state = StateClosed
			cb.failureCount = 0
			cb.successCount = 0
		}
	}

	cb.mu.Unlock()

	err := fn(ctx)

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		if cb.failureCount >= cb.threshold {
			cb.state = StateOpen
		}
		return err
	}

	cb.successCount++
	return nil
}

type Bulkhead struct {
	sem chan struct{}
}

func NewBulkhead(maxConcurrent int) *Bulkhead {
	return &Bulkhead{sem: make(chan struct{}, maxConcurrent)}
}

func (b *Bulkhead) Execute(ctx context.Context, fn func(context.Context) error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
		return fn(ctx)
	case <-ctx.Done():
		return ctx.Err()
	}
}

type DegradationFunc func(ctx context.Context) (string, error)

type ResilientClient struct {
	primary  DegradationFunc
	fallback DegradationFunc
	breaker  *CircuitBreaker
	bulkhead *Bulkhead
}

func NewResilientClient(primary, fallback DegradationFunc) *ResilientClient {
	return &ResilientClient{
		primary:  primary,
		fallback: fallback,
		breaker:  NewCircuitBreaker(5, 10*time.Second),
		bulkhead: NewBulkhead(10),
	}
}

func (c *ResilientClient) Fetch(ctx context.Context) (string, error) {
	var result string
	err := c.breaker.Call(ctx, func(ctx context.Context) error {
		return c.bulkhead.Execute(ctx, func(ctx context.Context) error {
			r, e := c.primary(ctx)
			result = r
			return e
		})
	})

	if err != nil {
		fallbackResult, fallbackErr := c.fallback(ctx)
		if fallbackErr != nil {
			return "", fmt.Errorf("all paths failed: %w", fallbackErr)
		}
		return fallbackResult + " (degraded)", nil
	}

	return result, nil
}

func main() {
	callCount := 0
	client := NewResilientClient(
		func(ctx context.Context) (string, error) {
			callCount++
			if callCount <= 3 {
				return "", errors.New("primary failure")
			}
			return "primary data", nil
		},
		func(ctx context.Context) (string, error) {
			return "cached data", nil
		},
	)

	fmt.Println("=== Resilient Client Demo ===")
	for i := 0; i < 8; i++ {
		result, err := client.Fetch(context.Background())
		state := client.breaker.State()
		if err != nil {
			fmt.Printf("Request %d: error [breaker=%s] — %v\n", i+1, state, err)
		} else {
			fmt.Printf("Request %d: %q [breaker=%s]\n", i+1, result, state)
		}
	}
}
