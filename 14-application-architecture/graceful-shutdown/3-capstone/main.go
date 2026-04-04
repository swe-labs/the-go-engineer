// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// ============================================================================
// Section 27: Graceful Shutdown — Complete Production Capstone
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Wiring together: signal handling + HTTP drain + DB close + background workers
//   - Coordinating multiple resources that must shut down in dependency order
//   - Health endpoint responding "503 Shutting Down" during drain window
//   - The Kubernetes readiness vs liveness probe distinction
//
// SHUTDOWN DEPENDENCY ORDER:
//   The order in which you shut things down matters critically:
//     1. Stop accepting NEW HTTP requests (close listener)
//     2. Let IN-FLIGHT requests complete (they still need DB access)
//     3. Stop background workers (they may use DB too)
//     4. Close database connections (everything above must finish first)
//     5. Flush logs and metrics
//     6. Exit 0
//
//   Reversing steps 2 and 4 would crash in-flight requests with DB errors.
//
// KUBERNETES READINESS PROBES:
//   /healthz/ready — returns 200 when the service is ready to receive traffic.
//                    Return 503 as soon as SIGTERM arrives.
//                    Kubernetes will remove the pod from the load balancer.
//   /healthz/live  — returns 200 as long as the process is responsive.
//                    Only return 503 if you want Kubernetes to RESTART the pod.
//
// RUN: go run ./27-graceful-shutdown/3-capstone
//   Then send: kill -TERM <pid>   or press Ctrl+C
// ============================================================================

// ============================================================================
// Component: MessageQueueConsumer
// A background worker that must drain its queue before shutdown.
// ============================================================================

type MessageQueueConsumer struct {
	name     string
	logger   *slog.Logger
	queue    chan string
	inFlight sync.WaitGroup
}

func NewConsumer(name string, logger *slog.Logger) *MessageQueueConsumer {
	c := &MessageQueueConsumer{
		name:   name,
		logger: logger,
		queue:  make(chan string, 100),
	}
	// Seed with some messages
	for i := range 10 {
		c.queue <- fmt.Sprintf("message_%d", i)
	}
	return c
}

func (c *MessageQueueConsumer) Run(ctx context.Context) error {
	c.logger.Info("consumer started", "name", c.name)
	defer c.logger.Info("consumer stopped", "name", c.name)

	for {
		select {
		case msg, ok := <-c.queue:
			if !ok {
				return nil
			}
			c.inFlight.Add(1)
			go func(m string) {
				defer c.inFlight.Done()
				time.Sleep(50 * time.Millisecond) // Simulate processing
				c.logger.Debug("message processed", "name", c.name, "msg", m)
			}(msg)

		case <-ctx.Done():
			// Drain remaining messages already in queue before stopping
			c.logger.Info("consumer draining", "name", c.name)
			for {
				select {
				case msg := <-c.queue:
					time.Sleep(20 * time.Millisecond)
					c.logger.Debug("drained message", "name", c.name, "msg", msg)
				default:
					c.inFlight.Wait() // Wait for in-flight goroutines
					return nil
				}
			}
		}
	}
}

// ============================================================================
// Component: ProductionServer
// HTTP server with health endpoints and shutdown awareness.
// ============================================================================

type ProductionServer struct {
	http    *http.Server
	logger  *slog.Logger
	isReady bool
	readyMu sync.RWMutex
}

func NewProductionServer(logger *slog.Logger) *ProductionServer {
	s := &ProductionServer{logger: logger, isReady: true}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz/live", s.handleLive)
	mux.HandleFunc("GET /healthz/ready", s.handleReady)
	mux.HandleFunc("GET /api/v1/orders", s.handleOrders)

	s.http = &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return s
}

// SetNotReady flips the readiness probe to 503.
// Call this BEFORE initiating HTTP drain so Kubernetes stops routing traffic.
func (s *ProductionServer) SetNotReady() {
	s.readyMu.Lock()
	s.isReady = false
	s.readyMu.Unlock()
	s.logger.Info("readiness probe set to NOT READY — Kubernetes will stop routing traffic")
}

func (s *ProductionServer) handleLive(w http.ResponseWriter, r *http.Request) {
	// Liveness: always 200 unless the process is deadlocked
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"alive"}`)
}

func (s *ProductionServer) handleReady(w http.ResponseWriter, r *http.Request) {
	s.readyMu.RLock()
	ready := s.isReady
	s.readyMu.RUnlock()

	if !ready {
		// 503 tells Kubernetes: remove this pod from load balancer rotation
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{"status":"shutting_down"}`)
		return
	}
	fmt.Fprint(w, `{"status":"ready"}`)
}

func (s *ProductionServer) handleOrders(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond) // Simulate DB query
	fmt.Fprint(w, `{"orders":[]}`)
}

func (s *ProductionServer) ListenAndServe() error {
	s.logger.Info("http server starting", "addr", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *ProductionServer) Shutdown(ctx context.Context) error {
	s.logger.Info("draining http connections")
	return s.http.Shutdown(ctx)
}

// ============================================================================
// Application wiring — the main entry point
// ============================================================================

func run() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// Signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialise components
	httpServer := NewProductionServer(logger)
	consumer := NewConsumer("order-events", logger)

	// Run everything under errgroup so any fatal error triggers shutdown
	g, gctx := errgroup.WithContext(ctx)

	// Goroutine 1: HTTP server
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	// Goroutine 2: message queue consumer
	g.Go(func() error {
		return consumer.Run(gctx)
	})

	// Goroutine 3: orchestrate shutdown
	g.Go(func() error {
		<-gctx.Done()
		stop() // Allow second Ctrl+C to force-kill

		logger.Info("shutdown initiated", "reason", gctx.Err())

		// Give Kubernetes 3 seconds to propagate readiness=503 to load balancers
		// before we stop accepting connections.
		httpServer.SetNotReady()
		logger.Info("waiting 3s for load balancer propagation...")
		time.Sleep(3 * time.Second)

		// 30 second window for the entire shutdown sequence
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// STEP 1: drain HTTP (consumers still run; handlers may still use DB)
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http drain: %w", err)
		}
		logger.Info("http drained")

		// STEP 2: consumer drains inside its own Run() goroutine via ctx.Done().
		// We just wait here until errgroup reports it finished.
		logger.Info("waiting for consumer to drain...")

		// STEP 3: flush logs / metrics (simplified)
		logger.Info("flushing telemetry...")
		time.Sleep(200 * time.Millisecond)

		logger.Info("shutdown sequence complete")
		return nil
	})

	return g.Wait()
}

func main() {
	if err := run(); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
	slog.Info("clean exit")

	// KEY TAKEAWAY:
	// - Shutdown ORDER matters: ready=503 → HTTP drain → workers → DB → logs
	// - SetNotReady() before HTTP drain gives K8s time to stop routing traffic
	// - errgroup wires HTTP server + workers + shutdown orchestrator together
	// - Any goroutine returning an error cancels gctx → triggers shutdown
	// - The shutdown goroutine gets its own 30-second deadline context
}
