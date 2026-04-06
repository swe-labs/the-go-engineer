// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// ============================================================================
// Section 24: errgroup & sync.Pool — URL Health Checker (Exercise Solution)
// Level: Advanced
// ============================================================================
//
// RUN: go run ./24-errgroup-and-pools/5-url-checker-exercise
// ============================================================================

// CheckResult holds the outcome of a single URL health check.
type CheckResult struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Error      error
}

// clientPool reuses HTTP clients to avoid allocating a new one per request.
// http.Client is safe for concurrent use — we pool it to avoid allocation,
// NOT to avoid data races.
var clientPool = sync.Pool{
	New: func() any {
		return &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 10,
			},
		}
	},
}

// checkURL performs a single HTTP HEAD request to url and returns the result.
// It uses the pooled client and respects context cancellation.
func checkURL(ctx context.Context, url string) CheckResult {
	start := time.Now()

	client := clientPool.Get().(*http.Client)
	defer clientPool.Put(client)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return CheckResult{URL: url, Error: err, Latency: time.Since(start)}
	}

	resp, err := client.Do(req)
	if err != nil {
		return CheckResult{URL: url, Error: err, Latency: time.Since(start)}
	}
	resp.Body.Close()

	return CheckResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		Latency:    time.Since(start),
	}
}

var urls = []string{
	"https://go.dev",
	"https://pkg.go.dev",
	"https://github.com",
	"https://api.github.com",
	"https://httpbin.org/status/200",
}

func main() {
	fmt.Println("=== URL Health Checker ===")
	fmt.Println()

	start := time.Now()
	results := make([]CheckResult, len(urls))

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(5) // Maximum 5 concurrent requests

	for i, url := range urls {
		i, url := i, url
		g.Go(func() error {
			results[i] = checkURL(ctx, url)
			return nil // We collect errors in CheckResult, not via errgroup
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("errgroup error: %v\n", err)
		return
	}

	// Sort by latency (fastest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Latency < results[j].Latency
	})

	fmt.Printf("%-45s %-8s %s\n", "URL", "STATUS", "LATENCY")
	fmt.Printf("%-45s %-8s %s\n", "---", "------", "-------")

	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("%-45s %-8s %s\n", r.URL, "ERROR", r.Error)
			continue
		}
		status := fmt.Sprintf("%d", r.StatusCode)
		icon := "✅"
		if r.StatusCode >= 400 {
			icon = "❌"
		}
		fmt.Printf("%s %-43s %-8s %v\n", icon, r.URL, status, r.Latency.Round(time.Millisecond))
	}

	fmt.Printf("\nTotal time: %v (would be %v sequential)\n",
		time.Since(start).Round(time.Millisecond),
		sumLatencies(results).Round(time.Millisecond))
}

func sumLatencies(results []CheckResult) time.Duration {
	var total time.Duration
	for _, r := range results {
		total += r.Latency
	}
	return total
}
