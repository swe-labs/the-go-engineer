package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func retryWithBackoff(ctx context.Context, cfg RetryConfig, fn func(context.Context) error) error {
	var err error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err = fn(ctx)
		if err == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts {
			return fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, err)
		}
		delay := cfg.BaseDelay * (1 << (attempt - 1))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay = delay/2 + jitter

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return err
}

func unstableOp(ctx context.Context) error {
	// Simulate failure 70% of the time
	if rand.Intn(100) < 70 {
		return fmt.Errorf("transient error")
	}
	return nil
}

func main() {
	ctx := context.Background()
	cfg := RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    1 * time.Second,
	}

	err := retryWithBackoff(ctx, cfg, unstableOp)
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
		return
	}
	fmt.Println("Operation succeeded")
}
