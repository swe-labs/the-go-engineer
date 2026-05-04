// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/events"
	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/workers"
)

// setupGracefulShutdown configures a background goroutine that listens for
// OS termination signals (SIGINT, SIGTERM) and coordinates a safe drain of
// all in-flight work before allowing the process to exit.
//
// The returned channel is closed when the shutdown sequence completes,
// signaling to the main thread that it is safe to close final resources
// (like the database) and exit.
func setupGracefulShutdown(
	server *http.Server,
	shutdownTimeout time.Duration,
	logger *slog.Logger,
	isDraining *atomic.Bool,
	bus *events.Bus,
	cancelApp context.CancelFunc,
	workerPools ...*workers.Pool,
) <-chan struct{} {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	return drainOnSignal(sigint, server, shutdownTimeout, logger, isDraining, bus, cancelApp, workerPools...)
}

func drainOnSignal(
	sigint <-chan os.Signal,
	server *http.Server,
	shutdownTimeout time.Duration,
	logger *slog.Logger,
	isDraining *atomic.Bool,
	bus *events.Bus,
	cancelApp context.CancelFunc,
	workerPools ...*workers.Pool,
) <-chan struct{} {
	idleConnsClosed := make(chan struct{})

	go func() {
		<-sigint

		logger.Info("shutdown signal received, initiating graceful drain")

		// 1. Mark the application as draining. The /health endpoint will now
		// return 503, telling the load balancer to stop sending new traffic.
		isDraining.Store(true)

		// 2. Shut down the HTTP server. This stops accepting new connections
		// and waits for in-flight requests to finish, up to the deadline.
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("HTTP server shutdown failed or timed out", slog.Any("error", err))
		} else {
			logger.Info("HTTP server successfully drained")
		}

		// 3. Close the event bus. This prevents any further async jobs
		// from being published by late-finishing HTTP handlers.
		if bus != nil {
			bus.Close()
			logger.Info("event bus closed to new publications")
		}

		// 4. Cancel the root application context, signaling all background
		// systems that a graceful exit is underway.
		if cancelApp != nil {
			cancelApp()
		}

		// 5. Stop all background worker pools. This signals workers to stop
		// accepting new jobs, drain their internal buffers, and exit.
		// Wait blocks until all workers in the pool have returned.
		for _, pool := range workerPools {
			if pool != nil {
				pool.Stop()
				logger.Info("worker pool drained and stopped", slog.String("pool", pool.Name()))
			}
		}

		// 6. Signal the main goroutine that the shutdown sequence is complete.
		close(idleConnsClosed)
	}()

	return idleConnsClosed
}
