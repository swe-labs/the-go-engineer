// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Health and Readiness Probes
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The difference between Liveness and Readiness probes.
//   - How to implement '/healthz' and '/readyz' endpoints.
//   - How to simulate dependency checks in your probes.
//
// WHY THIS MATTERS:
//   - Orchestrators like Kubernetes and AWS ECS rely on these probes to
//     know when to restart your server or start sending it traffic.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/9-health-and-readiness-probes
//
// KEY TAKEAWAY:
//   - Liveness means "I am alive". Readiness means "I am ready to work".
// ============================================================================

package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Stage 06: Backend - Health and Readiness Probes
//
//   - Liveness (/healthz): Should we restart this container?
//   - Readiness (/readyz): Should we send traffic to this container?
//   - Startup Probes: Is the app finished initializing?
//
// ENGINEERING DEPTH:
//   Probes are the "Pulse" of your application. A readiness probe should
//   check all critical dependencies (Database, Redis, Upstream APIs).
//   If one is down, the app should mark itself as "not ready" so the
//   load balancer stops sending it traffic, giving the dependency time
//   to recover without causing user-facing errors.

func main() {
	fmt.Println("=== Health and Readiness Probes ===")
	fmt.Println()

	// isReady is an atomic boolean to safely toggle readiness state.
	var isReady atomic.Bool
	isReady.Store(false)

	// Simulate a slow startup process (e.g., loading cache or connecting to DB).
	go func() {
		fmt.Println("  [STARTUP] App is initializing...")
		time.Sleep(5 * time.Second)
		isReady.Store(true)
		fmt.Println("  [STARTUP] App is now READY to handle traffic.")
	}()

	mux := http.NewServeMux()

	// 1. Liveness Probe (/healthz)
	// Usually very simple. If the server can respond at all, it's alive.
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 2. Readiness Probe (/readyz)
	// Checks if the app is ready to serve traffic.
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if isReady.Load() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ready"))
		} else {
			// 503 tells the load balancer NOT to send traffic here.
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Not Ready"))
		}
	})

	// 3. Admin endpoint to toggle readiness for demonstration
	mux.HandleFunc("POST /admin/toggle-ready", func(w http.ResponseWriter, r *http.Request) {
		current := isReady.Load()
		isReady.Store(!current)
		fmt.Fprintf(w, "Readiness toggled to: %v\n", !current)
	})

	fmt.Println("  Server starting on :8088...")
	fmt.Println("  1. Check liveness:  curl -i http://localhost:8088/healthz")
	fmt.Println("  2. Check readiness: curl -i http://localhost:8088/readyz")
	fmt.Println("  3. Toggle state:    curl -X POST http://localhost:8088/admin/toggle-ready")
	fmt.Println()

	err := http.ListenAndServe(":8088", mux)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.10 rest-api-exercise")
	fmt.Println("Current: HS.9 (health-and-readiness-probes)")
	fmt.Println("Previous: HS.8 (graceful-http-shutdown)")
	fmt.Println("---------------------------------------------------")
}
