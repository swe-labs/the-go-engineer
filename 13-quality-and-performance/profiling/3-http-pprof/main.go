// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	_ "net/http/pprof" // Side-effect import: registers /debug/pprof/* handlers
	"sync"
	"time"
)

// ============================================================================
// Section 25: Profiling — Live pprof HTTP Endpoint
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - net/http/pprof: exposing live profiling endpoints on a running server
//   - How to take a CPU profile from a production-like server under load
//   - Security considerations: never expose pprof on a public-facing port
//   - The goroutine and mutex profiles (often overlooked)
//
// AFTER RUNNING THIS PROGRAM:
//
//   # 5-second CPU profile (captures what the server is doing right now)
//   go tool pprof http://localhost:8080/debug/pprof/profile?seconds=5
//
//   # Heap / memory profile
//   go tool pprof http://localhost:8080/debug/pprof/heap
//
//   # All currently running goroutines (detect goroutine leaks)
//   curl http://localhost:8080/debug/pprof/goroutine?debug=2
//
//   # Mutex contention profile (see which mutexes are blocking goroutines)
//   go tool pprof http://localhost:8080/debug/pprof/mutex
//
//   # Open interactive web UI for any profile:
//   go tool pprof -http=:8090 http://localhost:8080/debug/pprof/profile?seconds=5
//
// SECURITY NOTE:
//   NEVER expose /debug/pprof on a public internet-facing port. pprof exposes
//   symbol names, goroutine stacks (which may contain secrets), and heap contents.
//   The standard pattern is to run a SEPARATE internal admin server on an
//   internal port that only VPN-connected engineers can reach:
//
//     go http.ListenAndServe(":6060", nil) // pprof on internal port
//     http.ListenAndServe(":8080", mux)    // public API on separate mux
//
// RUN: go run ./25-profiling/3-http-pprof
// ============================================================================

// simulateWork does CPU-intensive computation to make profiles interesting.
func simulateWork(intensity int) {
	var sum float64
	for i := 0; i < intensity*1000; i++ {
		sum += rand.Float64() * float64(i)
	}
	_ = sum
}

// contentionHotspot demonstrates mutex contention — visible in pprof mutex profile.
var sharedCache struct {
	sync.RWMutex
	data map[string]string
}

func initCache() {
	sharedCache.data = make(map[string]string, 1000)
	for i := 0; i < 1000; i++ {
		sharedCache.data[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}
}

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	// Simulate a compute-intensive operation (shows up in CPU profile)
	simulateWork(10)

	// Read from the shared cache (holds RLock briefly)
	sharedCache.RLock()
	val := sharedCache.data[fmt.Sprintf("key_%d", rand.Intn(1000))]
	sharedCache.RUnlock()

	// Occasionally write to the cache (holds full Lock — contention point)
	if rand.Intn(20) == 0 {
		key := fmt.Sprintf("key_%d", rand.Intn(1000))
		sharedCache.Lock()
		sharedCache.data[key] = fmt.Sprintf("updated_%d", time.Now().UnixNano())
		sharedCache.Unlock()
	}

	fmt.Fprintf(w, `{"status":"ok","result":"%s"}`, val)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"healthy"}`)
}

func generateLoad(duration time.Duration) {
	// Send a constant stream of requests to the server to make profiles non-trivial.
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.Now().Add(duration)
	var wg sync.WaitGroup

	for time.Now().Before(deadline) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Get("http://localhost:8080/api/v1/data")
			if err == nil {
				resp.Body.Close()
			}
		}()
		time.Sleep(5 * time.Millisecond) // 200 req/s
	}
	wg.Wait()
}

func main() {
	initCache()

	// =========================================================================
	// Two-port pattern: public API on :8080, pprof on :6060
	// =========================================================================
	// This is the production-safe way to run pprof. Port 6060 is only
	// accessible via VPN or kubectl port-forward — never exposed publicly.

	// Public API mux
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /api/v1/data", handleAPIRequest)
	apiMux.HandleFunc("GET /healthz", handleHealthz)

	// Internal admin mux — pprof handlers registered via the blank import above.
	// The blank import `_ "net/http/pprof"` registers all /debug/pprof/* handlers
	// on the DEFAULT ServeMux (http.DefaultServeMux). We serve them separately.
	// DO NOT add pprof to apiMux — that would expose it publicly.

	slog.Info("public API starting", "addr", ":8080")
	slog.Info("pprof admin starting", "addr", ":6060")
	slog.Info("take a CPU profile with: go tool pprof http://localhost:6060/debug/pprof/profile?seconds=5")

	// Start pprof server (internal only, uses DefaultServeMux)
	go func() {
		log.Fatal(http.ListenAndServe(":6060", nil)) // nil = DefaultServeMux with pprof
	}()

	// Generate load after a short delay to let the server start
	go func() {
		time.Sleep(500 * time.Millisecond)
		slog.Info("generating load for 30 seconds...")
		generateLoad(30 * time.Second)
		slog.Info("load generation complete")
	}()

	// Start public API server
	log.Fatal(http.ListenAndServe(":8080", apiMux))

	// KEY TAKEAWAY:
	// - `_ "net/http/pprof"` registers /debug/pprof/* on DefaultServeMux
	// - NEVER mix pprof with your public API mux — use a separate port
	// - go tool pprof -http=:8090 <url> gives the full interactive web UI
	// - Goroutine profile detects leaks: /debug/pprof/goroutine?debug=2
	// - Mutex profile shows lock contention: /debug/pprof/mutex
	// - Heap profile shows live allocations: /debug/pprof/heap
}
