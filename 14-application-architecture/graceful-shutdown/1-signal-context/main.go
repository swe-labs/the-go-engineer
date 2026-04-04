// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ============================================================================
// Section 27: Graceful Shutdown — signal.NotifyContext
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How OS signals work: SIGTERM (graceful stop) vs SIGKILL (force kill)
//   - signal.NotifyContext: the idiomatic Go 1.16+ signal handling API
//   - Why every production binary must handle SIGTERM
//   - The difference between the old signal.Notify pattern and NotifyContext
//
// OS SIGNALS — what you must know:
//   SIGINT  (Ctrl+C)    — sent when user presses Ctrl+C in terminal
//   SIGTERM             — sent by Kubernetes, systemd, docker stop, kill <pid>
//   SIGKILL             — cannot be caught. Immediate termination (last resort).
//
//   The deployment sequence:
//     1. docker stop / kubectl delete pod → sends SIGTERM
//     2. Process has 30 seconds to clean up
//     3. If still running after 30s → SIGKILL (force kill, data loss possible)
//
// ENGINEERING DEPTH:
//   Before signal.NotifyContext (Go 1.16), signal handling required:
//     sigCh := make(chan os.Signal, 1)
//     signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
//     go func() {
//         <-sigCh
//         // cleanup...
//     }()
//   This pattern has subtle bugs: signal.Notify(nil chan) panics, the goroutine
//   may not run before the signal arrives, and cleanup errors are ignored.
//
//   signal.NotifyContext wraps all of this in a context that is cancelled when
//   the signal arrives. This integrates naturally with every context-aware API
//   (database queries, HTTP clients, gRPC calls) — they all cancel automatically.
//
// RUN: go run ./27-graceful-shutdown/1-signal-context
//   Then press Ctrl+C to see the signal handling in action.
//   Or send: kill -TERM <pid>
// ============================================================================

// BackgroundWorker simulates a long-running background task (metrics exporter,
// message queue consumer, cache refresher, etc.)
func BackgroundWorker(ctx context.Context, name string) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	defer slog.Info("worker stopped", "name", name)

	slog.Info("worker started", "name", name)

	for {
		select {
		case <-ctx.Done():
			// Context cancelled — perform cleanup before returning
			slog.Info("worker received shutdown signal", "name", name, "reason", ctx.Err())
			// Simulate cleanup (flush buffer, close connection, etc.)
			time.Sleep(200 * time.Millisecond)
			return
		case t := <-ticker.C:
			slog.Debug("worker tick", "name", name, "time", t.Format("15:04:05"))
		}
	}
}

// Database simulates a database connection that must be closed gracefully.
type Database struct{ name string }

func (db *Database) Close() error {
	slog.Info("closing database connection", "db", db.name)
	time.Sleep(100 * time.Millisecond) // Simulate connection drain
	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// =========================================================================
	// signal.NotifyContext — the idiomatic modern pattern
	// =========================================================================
	// ctx is cancelled when SIGINT (Ctrl+C) or SIGTERM is received.
	// stop() stops the signal relay — call it as soon as the context is done
	// to allow a second Ctrl+C to force-kill if cleanup hangs.
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // SIGINT (Ctrl+C)
		syscall.SIGTERM, // sent by Kubernetes, docker stop, systemd
	)
	defer stop() // Release resources if we return before a signal

	// Connect to the database (cleanup must happen on shutdown)
	db := &Database{name: "postgres-primary"}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("database close error", "error", err)
		}
	}()

	// =========================================================================
	// Launch background workers using the signal context
	// =========================================================================
	// When SIGTERM arrives, ctx is cancelled, and ALL workers stop automatically.
	// No manual plumbing required — the context tree propagates the signal.
	go BackgroundWorker(ctx, "metrics-exporter")
	go BackgroundWorker(ctx, "cache-refresher")
	go BackgroundWorker(ctx, "health-broadcaster")

	slog.Info("application started — press Ctrl+C or send SIGTERM to stop")

	// =========================================================================
	// Block until signal arrives
	// =========================================================================
	<-ctx.Done()
	stop() // Call stop() early so a second Ctrl+C force-kills immediately

	slog.Info("shutdown signal received", "reason", ctx.Err())

	// =========================================================================
	// Graceful cleanup with a shutdown deadline
	// =========================================================================
	// Give background goroutines time to finish in-progress work.
	// After the deadline, exit anyway — don't wait forever.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Wait for all workers to finish (simplified — in production use errgroup or WaitGroup)
	done := make(chan struct{})
	go func() {
		// Simulate waiting for all workers
		time.Sleep(500 * time.Millisecond) // Workers need ~200-300ms to clean up
		close(done)
	}()

	select {
	case <-done:
		slog.Info("graceful shutdown complete")
	case <-shutdownCtx.Done():
		slog.Warn("shutdown deadline exceeded — forcing exit")
		fmt.Fprintln(os.Stderr, "shutdown timeout: some workers may not have finished cleanly")
		os.Exit(1)
	}

	// KEY TAKEAWAY:
	// - signal.NotifyContext: ctx cancelled on SIGINT or SIGTERM
	// - All context-aware goroutines stop automatically when ctx.Done() closes
	// - Call stop() after ctx.Done() so a 2nd Ctrl+C force-kills the process
	// - Use a separate context with a deadline for the cleanup phase itself
	// - defer db.Close() runs even when the process receives a signal
}
