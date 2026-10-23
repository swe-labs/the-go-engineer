package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestProbeTarget_SuccessAndFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthy" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("error"))
		}
	}))
	defer server.Close()

	checker := NewHealthChecker(CheckerConfig{
		Workers: 2,
		Timeout: 1 * time.Second,
		Client:  server.Client(),
	})

	ctx := t.Context()

	// Healthy
	resUp := checker.ProbeTarget(ctx, server.URL+"/healthy")
	if resUp.Status != StatusUp {
		t.Fatalf("expected StatusUp, got %s: %s", resUp.Status, resUp.Error)
	}
	if resUp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resUp.StatusCode)
	}

	// Unhealthy
	resDown := checker.ProbeTarget(ctx, server.URL+"/fail")
	if resDown.Status != StatusDown {
		t.Fatalf("expected StatusDown, got %s", resDown.Status)
	}
	if resDown.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resDown.StatusCode)
	}
}

func TestProbeTarget_TimeoutHandling(t *testing.T) {
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer slowServer.Close()

	// 50ms timeout < 200ms sleep
	checker := NewHealthChecker(CheckerConfig{
		Workers: 1,
		Timeout: 50 * time.Millisecond,
		Client:  slowServer.Client(),
	})

	ctx := t.Context()
	res := checker.ProbeTarget(ctx, slowServer.URL)

	if res.Status != StatusDown {
		t.Fatalf("expected StatusDown due to timeout, got %s", res.Status)
	}
	if res.Error == "" {
		t.Fatal("expected error message describing timeout")
	}
}

func TestCheckAll_ConcurrentExecution(t *testing.T) {
	var concurrentRequests atomic.Int32
	var peakConcurrent atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		curr := concurrentRequests.Add(1)
		defer concurrentRequests.Add(-1)

		for {
			peak := peakConcurrent.Load()
			if curr <= peak || peakConcurrent.CompareAndSwap(peak, curr) {
				break
			}
		}

		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	const workerCount = 4
	checker := NewHealthChecker(CheckerConfig{
		Workers: workerCount,
		Timeout: 2 * time.Second,
		Client:  server.Client(),
	})

	urls := make([]string, 16)
	for i := range urls {
		urls[i] = server.URL
	}

	results := checker.CheckAll(t.Context(), urls)

	if len(results) != len(urls) {
		t.Fatalf("expected %d results, got %d", len(urls), len(results))
	}

	for idx, r := range results {
		if r.Status != StatusUp {
			t.Fatalf("result %d failed: %s", idx, r.Error)
		}
	}

	if peakConcurrent.Load() > workerCount {
		t.Fatalf("peak concurrent requests (%d) exceeded configured workers (%d)", peakConcurrent.Load(), workerCount)
	}
}

func TestStatusTransitions(t *testing.T) {
	isHealthy := true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHealthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer server.Close()

	checker := NewHealthChecker(CheckerConfig{
		Workers: 1,
		Timeout: 1 * time.Second,
		Client:  server.Client(),
	})

	ctx := t.Context()

	// Initial probe: UP
	res1 := checker.ProbeTarget(ctx, server.URL)
	if res1.Status != StatusUp {
		t.Fatalf("expected initial probe UP, got %s", res1.Status)
	}

	// Change state to DOWN
	isHealthy = false
	res2 := checker.ProbeTarget(ctx, server.URL)
	if res2.Status != StatusDown {
		t.Fatalf("expected second probe DOWN, got %s", res2.Status)
	}

	// Change state back to UP
	isHealthy = true
	res3 := checker.ProbeTarget(ctx, server.URL)
	if res3.Status != StatusUp {
		t.Fatalf("expected third probe UP, got %s", res3.Status)
	}

	transitions := checker.Transitions()
	if len(transitions) != 2 {
		t.Fatalf("expected 2 status transitions, got %d", len(transitions))
	}

	// Transition 1: UP -> DOWN
	if transitions[0].FromStatus != StatusUp || transitions[0].ToStatus != StatusDown {
		t.Fatalf("expected transition 1 UP -> DOWN, got %s -> %s", transitions[0].FromStatus, transitions[0].ToStatus)
	}

	// Transition 2: DOWN -> UP
	if transitions[1].FromStatus != StatusDown || transitions[1].ToStatus != StatusUp {
		t.Fatalf("expected transition 2 DOWN -> UP, got %s -> %s", transitions[1].FromStatus, transitions[1].ToStatus)
	}
}
