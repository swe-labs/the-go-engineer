package main

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestRetryWithBackoff(t *testing.T) {
	tests := []struct {
		name        string
		maxAttempts int
		fn          func(context.Context) error
		wantErr     bool
	}{
		{
			name:        "succeeds on first try",
			maxAttempts: 3,
			fn:          func(ctx context.Context) error { return nil },
			wantErr:     false,
		},
		{
			name:        "succeeds after retry",
			maxAttempts: 5,
			fn: func(ctx context.Context) error {
				if rand.Intn(100) < 50 {
					return errors.New("transient")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name:        "fails all attempts",
			maxAttempts: 2,
			fn:          func(ctx context.Context) error { return errors.New("permanent") },
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := RetryConfig{
				MaxAttempts: tt.maxAttempts,
				BaseDelay:   time.Millisecond,
				MaxDelay:    10 * time.Millisecond,
			}
			err := retryWithBackoff(context.Background(), cfg, tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("retryWithBackoff() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestRetryContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   time.Millisecond,
		MaxDelay:    time.Millisecond,
	}
	err := retryWithBackoff(ctx, cfg, func(ctx context.Context) error {
		return errors.New("transient")
	})
	if err == nil {
		t.Error("expected context cancellation error")
	}
}
