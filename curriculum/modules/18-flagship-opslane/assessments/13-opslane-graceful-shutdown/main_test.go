package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestGracefulShutdown_InFlightCompleteAndNewRejected(t *testing.T) {
	// Core practice task requirement:
	// Send an interrupt/shutdown signal, verify that in-flight requests complete,
	// and confirm that new requests are rejected.

	// 1. Allocate an available local TCP port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to bind listener: %v", err)
	}
	addr := listener.Addr().String()

	coordinator := NewShutdownCoordinator(addr)

	// Start server on listener
	go func() {
		_ = coordinator.server.Serve(listener)
	}()

	client := &http.Client{Timeout: 3 * time.Second}

	// Verify server is responsive initially
	resPing, err := client.Get("http://" + addr + "/ping")
	if err != nil || resPing.StatusCode != http.StatusOK {
		t.Fatalf("server not healthy initially: %v", err)
	}
	_ = resPing.Body.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	var inFlightStatus int
	var inFlightBody string
	var inFlightErr error

	started := make(chan struct{})
	var startOnce sync.Once
	coordinator.onSlowTaskStart = func() {
		startOnce.Do(func() {
			close(started)
		})
	}

	// 2. Launch long in-flight request
	go func() {
		defer wg.Done()
		resp, err := client.Get("http://" + addr + "/slow-task")
		if err != nil {
			inFlightErr = err
			return
		}
		defer resp.Body.Close()
		inFlightStatus = resp.StatusCode
		b, _ := io.ReadAll(resp.Body)
		inFlightBody = string(b)
	}()

	// Wait deterministically until the request has entered handleSlowTask
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for slow-task to start")
	}

	// 3. Initiate Graceful Shutdown concurrently while request is in flight
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := coordinator.GracefulShutdown(shutdownCtx); err != nil {
		t.Fatalf("GracefulShutdown failed: %v", err)
	}

	// 4. In-flight request must finish successfully (200 OK)
	wg.Wait()
	if inFlightErr != nil {
		t.Fatalf("in-flight request failed prematurely: %v", inFlightErr)
	}
	if inFlightStatus != http.StatusOK || inFlightBody != "SLOW_TASK_COMPLETED" {
		t.Fatalf("expected in-flight request to complete with 200 OK, got status %d, body: %q", inFlightStatus, inFlightBody)
	}

	// 5. Subsequent new request after shutdown MUST be rejected
	resLate, err := client.Get("http://" + addr + "/ping")
	if err == nil {
		defer resLate.Body.Close()
		if resLate.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("expected new request to be rejected with 503 Service Unavailable, got %d", resLate.StatusCode)
		}
	}
	// Note: Connection refused is also a valid rejection outcome once listener terminates.
}
