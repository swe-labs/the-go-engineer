// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// ============================================================================
// Section 27: Graceful Shutdown — Production HTTP Server
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - http.Server.Shutdown(): drain in-flight requests without cutting them
//   - The complete production graceful shutdown pattern:
//       signal → stop accepting → drain requests → close DB → exit
//   - Why log.Fatal(http.ListenAndServe(...)) is wrong for production
//   - Testing graceful shutdown: how to verify requests aren't dropped
//
// THE DEPLOYMENT PROBLEM:
//   Rolling deployments update pods one at a time. When Kubernetes terminates
//   the old pod it:
//     1. Removes the pod from the load balancer (takes 2-5 seconds to propagate)
//     2. Immediately sends SIGTERM to the process
//
//   During that 2-5 second window, traffic is STILL arriving at the pod.
//   Without Shutdown(), those requests are instantly terminated → 502 errors.
//   With Shutdown(), the server waits for those requests to complete → 0 errors.
//
// ENGINEERING DEPTH:
//   http.Server.Shutdown(ctx) does four things:
//     1. Closes the listener: no new connections accepted
//     2. Closes idle connections
//     3. Waits for active connections to become idle (request completes)
//     4. Returns when all connections are idle or ctx expires
//
//   http.Server.ListenAndServe() returns http.ErrServerClosed when Shutdown
//   is called. This is the expected "clean" exit — not an error. Always
//   check for this: if err != nil && !errors.Is(err, http.ErrServerClosed)
//
// RUN: go run ./27-graceful-shutdown/2-http-server
//   Then: curl http://localhost:8080/api/slow   (simulates a slow 3s request)
//   While the request is in-flight, press Ctrl+C — graceful shutdown waits for it.
// ============================================================================

type Server struct {
	httpServer *http.Server
	db         *sql.DB // In real use: an actual *sql.DB
	logger     *slog.Logger
}

func NewServer(logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	s := &Server{logger: logger}

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/orders", s.handleListOrders)
	mux.HandleFunc("GET /api/slow", s.handleSlow)

	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// Always set timeouts on http.Server — defaults are infinite.
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return s
}

// Start runs the HTTP server and blocks until it is shut down.
// It returns http.ErrServerClosed on graceful shutdown — callers must handle this.
func (s *Server) Start() error {
	s.logger.Info("server starting", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server error: %w", err)
	}
	return nil
}

// Shutdown gracefully drains the server and closes resources.
// The ctx deadline is the MAXIMUM time we wait for in-flight requests.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("initiating graceful shutdown")

	// Step 1: Stop accepting new connections + drain existing ones
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}
	s.logger.Info("http server drained")

	// Step 2: Close database (after HTTP — DB queries are used by handlers)
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("db close: %w", err)
		}
		s.logger.Info("database connection closed")
	}

	return nil
}

// ============================================================================
// Handlers
// ============================================================================

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleListOrders(w http.ResponseWriter, r *http.Request) {
	// Simulate a quick database query
	time.Sleep(10 * time.Millisecond)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]any{
		{"id": "ord_001", "status": "shipped"},
		{"id": "ord_002", "status": "pending"},
	})
}

// handleSlow simulates a long-running request (3 seconds).
// Use this to test that graceful shutdown waits for it to complete.
func (s *Server) handleSlow(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("slow request started")
	select {
	case <-time.After(3 * time.Second):
		s.logger.Info("slow request completed")
		fmt.Fprint(w, `{"result":"completed"}`)
	case <-r.Context().Done():
		// Client disconnected or server is shutting down with a tight deadline
		s.logger.Warn("slow request cancelled", "reason", r.Context().Err())
		http.Error(w, "request cancelled", http.StatusServiceUnavailable)
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	srv := NewServer(logger)

	// =========================================================================
	// The complete graceful shutdown pattern using errgroup
	// =========================================================================
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gctx := errgroup.WithContext(ctx)

	// Goroutine 1: run the HTTP server
	g.Go(func() error {
		return srv.Start() // Blocks until Shutdown() is called
	})

	// Goroutine 2: wait for the signal, then shut down
	g.Go(func() error {
		<-gctx.Done() // Signal arrived OR the server goroutine errored out
		stop()        // Allow a 2nd Ctrl+C to force-kill

		// Give in-flight requests up to 30 seconds to complete.
		// In Kubernetes, terminationGracePeriodSeconds controls the SIGKILL deadline.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		return srv.Shutdown(shutdownCtx)
	})

	// Wait for both goroutines.
	// If the server fails (non-ErrServerClosed error), gctx is cancelled
	// and the shutdown goroutine runs immediately.
	if err := g.Wait(); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}

	logger.Info("graceful shutdown complete")

	// KEY TAKEAWAY:
	// - http.Server.Shutdown(ctx): drains in-flight requests cleanly
	// - ListenAndServe returns http.ErrServerClosed on graceful shutdown (not an error)
	// - Always set ReadTimeout, WriteTimeout, IdleTimeout on http.Server
	// - errgroup orchestrates server + signal handler goroutines elegantly
	// - Kubernetes terminationGracePeriodSeconds = your maximum shutdown window
}
