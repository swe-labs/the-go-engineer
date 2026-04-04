// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// ============================================================================
// Section 24: errgroup & sync.Pool — errgroup Basics
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - errgroup.Group: the idiomatic replacement for sync.WaitGroup when
//     goroutines can fail
//   - How errgroup collects errors without channels or mutexes
//   - The difference between WaitGroup and errgroup
//   - Go() vs TryGo(): bounded concurrency without a semaphore channel
//
// WHY errgroup BEATS WaitGroup + error channel:
//   WaitGroup pattern (verbose, error-prone):
//     var wg sync.WaitGroup
//     errc := make(chan error, N)      // manual buffering
//     for _, item := range items {
//         wg.Add(1)
//         go func(v Item) {
//             defer wg.Done()
//             if err := process(v); err != nil {
//                 errc <- err          // easy to forget
//             }
//         }(item)
//     }
//     wg.Wait()
//     close(errc)
//     for err := range errc { ... }  // collects ALL errors
//
//   errgroup pattern (idiomatic):
//     var g errgroup.Group
//     for _, item := range items {
//         g.Go(func() error { return process(item) })
//     }
//     if err := g.Wait(); err != nil { ... }  // first non-nil error
//
// ENGINEERING DEPTH:
//   errgroup.Group stores exactly one error: the first non-nil error returned
//   by any goroutine. All other errors are discarded. This is the right default
//   because returning ALL errors is rarely useful — a single failed database
//   connection explains the cascade of failures that follows it.
//
//   If you genuinely need ALL errors, use errors.Join() in each goroutine to
//   build a combined error, or collect into a []error with a mutex.
//
// RUN: go run ./24-errgroup-and-pools/1-errgroup
// ============================================================================

// Service represents a microservice dependency that must be checked at startup.
type Service struct {
	Name    string
	Healthy bool
	Latency time.Duration
}

// checkService simulates a health check that can fail.
func checkService(name string) (*Service, error) {
	// Simulate varying latency
	latency := time.Duration(len(name)*50) * time.Millisecond
	time.Sleep(latency)

	// Simulate "payment-gateway" being unhealthy
	if name == "payment-gateway" {
		return nil, fmt.Errorf("%s: connection refused (port 9443)", name)
	}

	return &Service{Name: name, Healthy: true, Latency: latency}, nil
}

func main() {
	// =========================================================================
	// Demo 1: Basic errgroup — all goroutines, first error wins
	// =========================================================================
	fmt.Println("=== errgroup: startup health checks ===")
	start := time.Now()

	services := []string{"postgres", "redis", "payment-gateway", "auth-service", "email-service"}
	results := make([]*Service, len(services))

	// errgroup.Group is a zero-value-ready struct.
	// Like WaitGroup, you do NOT need New() for basic use.
	var g errgroup.Group

	for i, svc := range services {
		i, svc := i, svc // Capture loop variable (see Section 09 closure bug)
		g.Go(func() error {
			result, err := checkService(svc)
			if err != nil {
				return err // errgroup stores the first non-nil error
			}
			results[i] = result // Safe: each goroutine writes a unique index
			return nil
		})
	}

	// Wait() blocks until all goroutines finish and returns the first error.
	if err := g.Wait(); err != nil {
		fmt.Printf("❌ Health check failed (%.0fms): %v\n", float64(time.Since(start).Milliseconds()), err)
	} else {
		fmt.Printf("✅ All services healthy (%.0fms)\n", float64(time.Since(start).Milliseconds()))
		for _, r := range results {
			if r != nil {
				fmt.Printf("   %-20s latency: %v\n", r.Name, r.Latency)
			}
		}
	}

	fmt.Println()

	// =========================================================================
	// Demo 2: errors.Join — collecting ALL errors when you need them
	// =========================================================================
	// errgroup only keeps the first error. If you need all of them, collect
	// them inside each goroutine and join into a single compound error.
	fmt.Println("=== Collecting all errors with errors.Join ===")

	var g2 errgroup.Group
	errCh := make(chan error, len(services)) // Buffered so goroutines don't block

	for _, svc := range services {
		svc := svc
		g2.Go(func() error {
			_, err := checkService(svc)
			if err != nil {
				errCh <- err // Send to channel; don't return — keep running
			}
			return nil // Return nil so other goroutines continue
		})
	}

	g2.Wait()
	close(errCh)

	var allErrs []error
	for err := range errCh {
		allErrs = append(allErrs, err)
	}

	if combined := errors.Join(allErrs...); combined != nil {
		fmt.Printf("Found %d failing services:\n%v\n", len(allErrs), combined)
	}

	fmt.Println()

	// =========================================================================
	// Demo 3: TryGo — bounded concurrency without a semaphore channel
	// =========================================================================
	// g.TryGo(f) tries to start f as a goroutine. It succeeds ONLY if the
	// current goroutine count is below the limit set by SetLimit().
	// This is cleaner than the chan struct{} semaphore pattern.
	fmt.Println("=== TryGo: bounded concurrency (max 2 at a time) ===")

	var g3 errgroup.Group
	g3.SetLimit(2) // Maximum 2 goroutines running simultaneously

	for i := range 5 {
		i := i
		// TryGo returns false if the limit is reached — spin until space opens.
		for !g3.TryGo(func() error {
			fmt.Printf("  worker %d started\n", i)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("  worker %d done\n", i)
			return nil
		}) {
			time.Sleep(5 * time.Millisecond) // Brief yield before retrying
		}
	}
	g3.Wait()

	// KEY TAKEAWAY:
	// - errgroup.Group is sync.WaitGroup + first-error collection
	// - g.Go(f) runs f in a goroutine; g.Wait() returns the first error
	// - Use errors.Join() in each goroutine if you need ALL errors
	// - g.SetLimit(N) + g.TryGo(f) = bounded concurrency without chan struct{}
	// - Use Go: github.com/rasel9t6/the-go-engineer already has x/sync
}
