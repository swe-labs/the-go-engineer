// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Basic GET Requests
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to make HTTP GET requests using 'http.Get'.
//   - How to stream-parse JSON responses directly from the response body.
//   - The hidden dangers of 'http.DefaultClient'.
//
// WHY THIS MATTERS:
//   - Most Go services need to communicate with external APIs. Understanding
//     how to do this efficiently and safely is a foundational skill.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/http-client/1-get-posts
//
// KEY TAKEAWAY:
//   - Never use 'http.Get' in production because it has no timeouts.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Stage 06: HTTP Client - Basic GET Requests
//
//   - Making HTTP GET requests with http.Get
//   - Parsing JSON responses directly from resp.Body (io.Reader)
//   - The hidden danger of http.DefaultClient
//
// ENGINEERING DEPTH:
//   `http.Get()` uses the `http.DefaultClient`.
//   WARNING: The DefaultClient has NO TIMEOUT!
//   If the server is hanging, your Goroutine will hang FOREVER.
//   In production, ALWAYS create a custom `http.Client` with a timeout.

type RemoteDevice struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

// fetchDevices fetches data from an external API.
// NOTE: This code is currently difficult to test because it hardcodes http.Get
// and a URL. We will refactor this in the next lesson.
func fetchDevices(limit int) ([]RemoteDevice, error) {
	// Using a public dummy API for demonstration
	url := fmt.Sprintf("https://dummyjson.com/products?limit=%d", limit)

	// DANGEROUS: Uses DefaultClient (no timeout).
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to reach external API: %w", err)
	}
	// ALWAYS defer closing the response body to prevent connection leaks!
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned unexpected status: %d", resp.StatusCode)
	}

	// Internal mapping for the dummy API response
	type DummyProduct struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Category string `json:"category"`
	}
	type DummyResponse struct {
		Products []DummyProduct `json:"products"`
	}

	var parsed DummyResponse

	// Use json.NewDecoder since resp.Body is an io.Reader (stream parsing).
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to decode JSON stream: %w", err)
	}

	var devices []RemoteDevice
	for _, p := range parsed.Products {
		devices = append(devices, RemoteDevice{
			ID:     p.ID,
			Type:   p.Category,
			Status: strings.ToUpper(p.Title),
		})
	}

	return devices, nil
}

func main() {
	fmt.Println("=== Fetching External API Data ===")
	fmt.Println()

	// Example of a safe client (how you should do it in production)
	_ = &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Println("  Initiating network request...")

	devices, err := fetchDevices(3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n  Successfully fetched %d devices:\n", len(devices))
	for _, d := range devices {
		fmt.Printf("  [ID: %3d] Type: %-15s Status: %s\n", d.ID, d.Type, d.Status)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  - Never use http.Get / DefaultClient in production (no timeouts)")
	fmt.Println("  - Always defer resp.Body.Close() to prevent socket leaks")
	fmt.Println("  - Use json.NewDecoder(resp.Body) for memory-efficient parsing")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.2 refactor-for-testability")
	fmt.Println("Current: HC.1 (get-posts)")
	fmt.Println("Previous: HS.10 (rest-api-exercise)")
	fmt.Println("---------------------------------------------------")
}
