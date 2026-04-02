// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// ============================================================================
// Section 16: HTTP Clients & Mocking — Basic GET Requests
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Making HTTP GET requests with `http.Get`
//   - Parsing JSON responses directly from `resp.Body` (io.Reader)
//   - The hidden danger of `http.DefaultClient`
//
// ANALOGY:
//   Using http.Get is like asking a stranger to fetch you a coffee.
//   They might come back in 5 minutes, or they might get hit by a bus
//   and never return, leaving you waiting forever.
//
// ENGINEERING DEPTH:
//   `http.Get()` uses the `http.DefaultClient`.
//   WARNING: The DefaultClient has NO TIMEOUT!
//   If the server is hanging, your Goroutine will hang FOREVER.
//   In production, ALWAYS create a custom `http.Client` with a timeout.
//
// RUN: go run ./16-http-clients-and-mocking/1-get-posts
// ============================================================================

type RemoteDevice struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type APIResponse struct {
	Devices []RemoteDevice `json:"devices"` // Matches the JSON "devices" array
}

// fetchDevices fetches data from an external API.
// NOTE: This code is currently UNTESTABLE because it hardcodes `http.Get`
// and hardcodes the URL. We cannot mock the network!
func fetchDevices(limit int) ([]RemoteDevice, error) {
	// Let's pretend this is a real external API returning devices
	url := fmt.Sprintf("https://dummyjson.com/products?limit=%d", limit)

	// ❌ DANGEROUS IN PRODUCTION: Uses DefaultClient (no timeout)
	// Hard dependency on the real internet.
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to reach external API: %w", err)
	}
	// ALWAYS defer closing the response body, or you will leak connections!
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned unexpected status: %d", resp.StatusCode)
	}

	// For educational purposes, since dummyjson returns "products", not "devices",
	// let's do a trick using arbitrary JSON mapping.
	type DummyProduct struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Category string `json:"category"`
	}
	type DummyResponse struct {
		Products []DummyProduct `json:"products"`
	}

	var parsed DummyResponse

	// Use json.NewDecoder since resp.Body is an io.Reader (stream parsing!)
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to decode JSON stream: %w", err)
	}

	// Map external API response to our domain model
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

	// Create a safe client demo (how you SHOULD do it in prod)
	_ = &http.Client{
		Timeout: 10 * time.Second, // Always set timeouts!
	}
	fmt.Println("   Warning: Our fetch logic currently uses DefaultClient (no timeout).")
	fmt.Println("   Initiating network request...")

	devices, err := fetchDevices(3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n   Successfully fetched %d devices:\n", len(devices))
	for _, d := range devices {
		fmt.Printf("   [ID: %3d] Type: %-20s Status: %s\n", d.ID, d.Type, d.Status)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Never use http.Get / DefaultClient in production (no timeouts).")
	fmt.Println("  - Always defer resp.Body.Close() to prevent socket leaks.")
	fmt.Println("  - Use json.NewDecoder(resp.Body) to stream-parse the response.")
	fmt.Println("  - Hardcoded http.Get calls make your function impossible to unit test!")
	fmt.Println("  - (Proceed to 2-refactor-for-testability to see the solution)")
}
