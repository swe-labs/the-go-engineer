// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Graceful HTTP Shutdown
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to handle OS signals (SIGINT, SIGTERM) to stop the server.
//   - How to use 'server.Shutdown' to let in-flight requests finish cleanly.
//   - How to set a deadline for the shutdown process.
//
// WHY THIS MATTERS:
//   - Simply killing a server process causes dropped connections and
//     potentially corrupted data. Graceful shutdown ensures a professional,
//     reliable user experience during deployments.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-servers/8-graceful-http-shutdown
//
// KEY TAKEAWAY:
//   - Close the listener first, then wait for workers to finish, then exit.
// ============================================================================

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Stage 06: Backend - Graceful HTTP Shutdown
//
//   - signal.Notify: Capturing Ctrl+C and kill signals
//   - server.Shutdown(ctx): The standard library cleanup method
//   - context.WithTimeout: Setting a hard limit on cleanup time
//
// ENGINEERING DEPTH:
//   When `server.Shutdown` is called, it immediately closes all active
//   listeners. It then waits for all active connections to become idle
//   and then closes them. This prevents new requests from starting while
//   giving existing ones time to finish their work. This is the difference
//   between a "Hard Crash" and a "Polite Exit".

func main() {
	fmt.Println("=== Graceful HTTP Shutdown ===")
	fmt.Println()

	mux := http.NewServeMux()

	// A handler that takes a few seconds to simulate real work.
	// Try hitting this endpoint and then immediately stopping the server.
	mux.HandleFunc("GET /work", func(w http.ResponseWriter, r *http.Request) {
		log.Println("  [WORK] Starting heavy work...")
		time.Sleep(5 * time.Second)
		fmt.Fprintln(w, "Work complete!")
		log.Println("  [WORK] Finished heavy work.")
	})

	server := &http.Server{
		Addr:    ":8087",
		Handler: mux,
	}

	// 1. Create a channel to listen for OS signals.
	// We use a buffered channel to ensure we don't miss any signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// 2. Start the server in a separate goroutine so it doesn't block.
	go func() {
		fmt.Println("  Server starting on :8087...")
		fmt.Println("  1. Visit: http://localhost:8087/work")
		fmt.Println("  2. Press Ctrl+C in this terminal immediately")
		fmt.Println("  3. Observe how the server waits for the work to finish!")
		fmt.Println()

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("  Critical error: %v", err)
		}
	}()

	// 3. Block here until we receive a signal.
	sig := <-stop
	fmt.Printf("\n  Received signal: %v. Starting graceful shutdown...\n", sig)

	// 4. Create a deadline for the shutdown process.
	// We give the server 10 seconds to finish all active requests.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5. Shutdown the server.
	// If the context expires before shutdown is complete, Shutdown returns
	// the context's error.
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("  Shutdown error: %v\n", err)
	} else {
		fmt.Println("  Server exited cleanly.")
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HS.9 -> 06-backend-db/01-web-and-database/http-servers/9-health-and-readiness-probes")
	fmt.Println("Current: HS.8 (graceful-http-shutdown)")
	fmt.Println("Previous: HS.7 (server-timeouts)")
	fmt.Println("---------------------------------------------------")
}
