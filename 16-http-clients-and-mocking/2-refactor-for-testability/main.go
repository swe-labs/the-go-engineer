// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// ============================================================================
// Section 16: HTTP Clients & Mocking — Refactor for Testability
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Dependency Injection for HTTP Clients
//   - Defining interfaces to mock external dependencies
//   - Structuring struct-based API Clients
//
// ANALOGY:
//   Instead of wiring our house directly to the town's power plant
//   (hardcoded `http.Get`), we install a standard wall outlet (interface).
//   Now we can plug in the real power grid (production build), OR we can
//   plug in a generator (a mock) during power outages (tests).
//
// ENGINEERING DEPTH:
//   Go implicitly satisfies interfaces. We can define our own interface
//   that exactly matches the signature of `(*http.Client).Do` or `(*http.Client).Get`.
//
// RUN: go run ./16-http-clients-and-mocking/2-refactor-for-testability
// ============================================================================

// --- 1. The Interface (The Wall Outlet) ---

// HTTPClient defines the exact signature of the `Get` method on `http.Client`.
// Because `*http.Client` already has this method, it implicitly satisfies this interface!
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// --- 2. The Client Struct ---

// IoTClient wraps our dependencies.
type IoTClient struct {
	// We depend on the INTERFACE, not the concrete `*http.Client`.
	// This means in our tests, we can provide a fake struct that implements Get().
	httpClient HTTPClient
	baseURL    string
}

// NewIoTClient acts as a constructor.
func NewIoTClient(client HTTPClient, baseURL string) *IoTClient {
	return &IoTClient{
		httpClient: client,
		baseURL:    baseURL,
	}
}

// --- 3. The Refactored Function ---

type SensorData struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

type APIResponse struct {
	Products []SensorData `json:"products"`
}

// FetchSensors now belongs to IoTClient. It uses c.httpClient, not http.Get!
func (c *IoTClient) FetchSensors(limit int) ([]SensorData, error) {
	// Construct URL dynamically so we can test different endpoints
	url := fmt.Sprintf("%s/products?limit=%d", c.baseURL, limit)

	// MAGIC: c.httpClient is just an interface.
	// In prod: hits the real internet. In tests: hits our mock.
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sensors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return apiResp.Products, nil
}

// --- 4. Main Execution ---

func main() {
	fmt.Println("=== Injecting HTTP Clients ===")

	baseURL := "https://dummyjson.com"

	// Create a safe, timeout-controlled HTTP client
	safeClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Inject the `safeClient` (which satisfies the HTTPClient interface)
	client := NewIoTClient(safeClient, baseURL)

	fmt.Println("   Fetching sensor data over the real network...")
	sensors, err := client.FetchSensors(3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error: %v\n", err)
		os.Exit(1)
	}

	for _, s := range sensors {
		fmt.Printf("   [Sensor %d] %s: $%.2f\n", s.ID, s.Title, s.Price)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Interface Injection decouples your code from the real network.")
	fmt.Println("  - Require an interface (HTTPClient); Pass a concrete struct (&http.Client).")
	fmt.Println("  - This refactor makes the code 100% unit-testable without mocking servers.")
}
