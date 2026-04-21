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
// Stage 07: Concurrency Patterns - errgroup Basics
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - errgroup.Group: the idiomatic replacement for sync.WaitGroup when goroutines can fail
//   - How errgroup collects errors without extra channels or mutexes
//   - The difference between WaitGroup and errgroup
//   - Go vs TryGo and SetLimit for bounded concurrency
//
// ENGINEERING DEPTH:
//   errgroup.Group stores exactly one error: the first non-nil error returned
//   by any goroutine. All other errors are discarded. This is the right default
//   because returning every error is rarely useful once one failed dependency
//   already explains the cascade of failures that follows it.
//
// RUN: go run ./07-concurrency/02-concurrency-patterns/1-errgroup
// ============================================================================

type Service struct {
	Name    string
	Healthy bool
	Latency time.Duration
}

func checkService(name string) (*Service, error) {
	latency := time.Duration(len(name)*50) * time.Millisecond
	time.Sleep(latency)

	if name == "payment-gateway" {
		return nil, fmt.Errorf("%s: connection refused (port 9443)", name)
	}

	return &Service{Name: name, Healthy: true, Latency: latency}, nil
}

func main() {
	fmt.Println("=== errgroup: startup health checks ===")
	start := time.Now()

	services := []string{"postgres", "redis", "payment-gateway", "auth-service", "email-service"}
	results := make([]*Service, len(services))

	var g errgroup.Group

	for i, svc := range services {
		i, svc := i, svc
		g.Go(func() error {
			result, err := checkService(svc)
			if err != nil {
				return err
			}
			results[i] = result
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("[FAIL] Health check failed (%dms): %v\n", time.Since(start).Milliseconds(), err)
	} else {
		fmt.Printf("[OK] All services healthy (%dms)\n", time.Since(start).Milliseconds())
		for _, r := range results {
			if r != nil {
				fmt.Printf("   %-20s latency: %v\n", r.Name, r.Latency)
			}
		}
	}

	fmt.Println()
	fmt.Println("=== Collecting all errors with errors.Join ===")

	var g2 errgroup.Group
	errCh := make(chan error, len(services))

	for _, svc := range services {
		svc := svc
		g2.Go(func() error {
			_, err := checkService(svc)
			if err != nil {
				errCh <- err
			}
			return nil
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
	fmt.Println("=== TryGo: bounded concurrency (max 2 at a time) ===")

	var g3 errgroup.Group
	g3.SetLimit(2)

	for i := range 5 {
		i := i
		for !g3.TryGo(func() error {
			fmt.Printf("  worker %d started\n", i)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("  worker %d done\n", i)
			return nil
		}) {
			time.Sleep(5 * time.Millisecond)
		}
	}
	g3.Wait()

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CP.2 errgroup + context")
	fmt.Println("   Current: CP.1 (errgroup basics)")
	fmt.Println("---------------------------------------------------")
}
