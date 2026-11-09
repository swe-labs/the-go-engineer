package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// ShutdownCoordinator coordinates graceful termination across HTTP and background resources.
type ShutdownCoordinator struct {
	server          *http.Server
	shuttingDown    atomic.Bool
	workerActive    atomic.Int32
	workerQueue     chan string
	cleanupDone     chan struct{}
	onSlowTaskStart func()
}

func NewShutdownCoordinator(addr string) *ShutdownCoordinator {
	coord := &ShutdownCoordinator{
		workerQueue: make(chan string, 100),
		cleanupDone: make(chan struct{}),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/slow-task", coord.handleSlowTask)
	mux.HandleFunc("/ping", coord.handlePing)

	coord.server = &http.Server{
		Addr:    addr,
		Handler: coord.healthMiddleware(mux),
	}

	return coord
}

func (sc *ShutdownCoordinator) healthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sc.shuttingDown.Load() {
			w.Header().Set("Connection", "close")
			http.Error(w, "503 Service Unavailable: Server is shutting down", http.StatusServiceUnavailable)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (sc *ShutdownCoordinator) handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("PONG"))
}

func (sc *ShutdownCoordinator) handleSlowTask(w http.ResponseWriter, r *http.Request) {
	if sc.onSlowTaskStart != nil {
		sc.onSlowTaskStart()
	}
	// Simulate an active in-flight request processing
	time.Sleep(50 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("SLOW_TASK_COMPLETED"))
}

// GracefulShutdown executes the phased shutdown sequence within a deadline context.
func (sc *ShutdownCoordinator) GracefulShutdown(ctx context.Context) error {
	// Step 1: Reject new requests at middleware level
	sc.shuttingDown.Store(true)

	// Step 2: Stop accepting TCP connections and wait for in-flight HTTP requests to drain
	if err := sc.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown failed: %w", err)
	}

	// Step 3: Flush background worker queue
	close(sc.workerQueue)

	return nil
}

func main() {
	coordinator := NewShutdownCoordinator(":8080")

	// Start server in background
	go func() {
		if err := coordinator.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server Listen error: %v", err)
		}
	}()

	fmt.Println("Opslane server listening on :8080. Press Ctrl+C to trigger graceful shutdown.")

	// Listen for OS interrupt signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	fmt.Println("\nSignal received. Draining in-flight requests...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := coordinator.GracefulShutdown(shutdownCtx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	fmt.Println("Graceful shutdown completed successfully.")
}
