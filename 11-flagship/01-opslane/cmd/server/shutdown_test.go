// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/events"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/workers"
)

func TestGracefulShutdownCoordination(t *testing.T) {
	// Setup dependencies
	server := &http.Server{}
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	isDraining := &atomic.Bool{}
	bus := events.NewBus(10)
	
	ctx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	pool, err := workers.NewPool(workers.PoolConfig{
		Name:      "test-pool",
		Workers:   1,
		QueueSize: 10,
		Handler: func(ctx context.Context, e events.Event) error {
			return nil
		},
	})
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	_ = pool.Start(ctx)

	sigint := make(chan os.Signal, 1)

	// Call drainOnSignal instead of setupGracefulShutdown to pass the sigint channel
	idleConnsClosed := drainOnSignal(sigint, server, 2*time.Second, logger, isDraining, bus, cancelApp, pool)

	// Verify pre-shutdown state
	if isDraining.Load() {
		t.Error("expected isDraining to be false initially")
	}

	// Trigger shutdown
	sigint <- os.Interrupt

	// Wait for shutdown to complete
	select {
	case <-idleConnsClosed:
		// success
	case <-time.After(3 * time.Second):
		t.Fatal("shutdown sequence timed out")
	}

	// 1. Verify draining flag
	if !isDraining.Load() {
		t.Error("expected isDraining to be true after shutdown")
	}

	// 2. Verify root context canceled
	if err := ctx.Err(); !errors.Is(err, context.Canceled) {
		t.Errorf("expected root context to be canceled, got %v", err)
	}

	// 3. Verify event bus closed
	err = bus.Publish(context.Background(), events.Event{Type: "test"})
	if !errors.Is(err, events.ErrBusClosed) {
		t.Errorf("expected event bus to be closed, got err: %v", err)
	}

	// 4. Verify worker pool stopped
	err = pool.Submit(context.Background(), events.Event{Type: "test"})
	if !errors.Is(err, workers.ErrPoolStopped) {
		t.Errorf("expected pool to be stopped, got err: %v", err)
	}
}
