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
// Stage 07: Concurrency Patterns - URL Health Checker (Exercise Solution)
// Level: Advanced
// ============================================================================
//
// RUN: go run ./07-concurrency/02-concurrency-patterns/5-url-checker-exercise
// ============================================================================

type CheckResult struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Error      error
}

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
	g.SetLimit(5)

	for i, url := range urls {
		i, url := i, url
		g.Go(func() error {
			results[i] = checkURL(ctx, url)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("errgroup error: %v\n", err)
		return
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Latency < results[j].Latency
	})

	fmt.Printf("%-7s %-45s %-8s %s\n", "RESULT", "URL", "STATUS", "LATENCY")
	fmt.Printf("%-7s %-45s %-8s %s\n", "------", "---", "------", "-------")

	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("%-7s %-45s %-8s %s\n", "ERROR", r.URL, "-", r.Error)
			continue
		}

		result := "OK"
		if r.StatusCode >= 400 {
			result = "FAIL"
		}
		fmt.Printf("%-7s %-45s %-8d %v\n", result, r.URL, r.StatusCode, r.Latency.Round(time.Millisecond))
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
