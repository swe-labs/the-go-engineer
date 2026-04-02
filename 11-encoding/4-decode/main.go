// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

// ============================================================================
// Section 11: Encoding — JSON Decoder (Streaming)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - json.NewDecoder: streaming JSON directly from an io.Reader
//   - Decoder vs Unmarshal: when to use which
//   - Decoding multiple JSON objects from a single stream (JSONL)
//   - Handling EOF (End Of File) properly
//
// ENGINEERING DEPTH:
//   A standard `json.Unmarshal` requires reading an entire API request body
//   into RAM before parsing can even begin. If a malicious user uploads a 5GB
//   JSON file, your container will Out-Of-Memory (OOM) crash instantly.
//   `json.NewDecoder` pulls bytes through the `io.Reader` interface lazily.
//   It maintains a tiny sliding memory buffer, parsing tokens on the fly and
//   discarding the bytes as it completes the struct layout. It is immune to
//   gigantic payloads.
//
// WHEN TO USE WHICH?
//   - Use json.Unmarshal if you already have the data in memory (a []byte)
//   - Use json.NewDecoder if you are reading from an io.Reader (like an HTTP Request Body or a File).
//     This avoids loading the entire payload into RAM before parsing.
//
// RUN: go run ./11-encoding/4-decode
// ============================================================================

type MetricEvent struct {
	AppID   string `json:"app_id"`
	Latency int    `json:"latency_ms"`
	Success bool   `json:"success"`
}

func main() {
	fmt.Println("=== JSON Decoder (Streaming) ===")
	fmt.Println()

	// =====================================================================
	// 1. Basic Decoder from an io.Reader
	// =====================================================================
	fmt.Println("1️⃣  Standard Decoder:")

	singlePayload := `{"app_id": "auth-service", "latency_ms": 125, "success": true}`

	// Convert string to io.Reader
	reader := strings.NewReader(singlePayload)

	// Create decoder
	dec := json.NewDecoder(reader)

	var event MetricEvent
	if err := dec.Decode(&event); err != nil {
		log.Fatal("Decode error:", err)
	}

	fmt.Printf("   Decoded: %+v\n\n", event)

	// =====================================================================
	// 2. Continuous Stream Parsing (JSON Lines / NDJSON)
	// =====================================================================
	// When reading from a network stream or log file, you might receive
	// multiple JSON objects separated by whitespace. json.NewDecoder handles
	// this gracefully, decoding one at a time.
	fmt.Println("2️⃣  Stream Parsing (Multiple Objects):")

	// Notice: these are three separate JSON objects, not a JSON array
	streamPayload := `
		{"app_id": "payment", "latency_ms": 45, "success": true}
		{"app_id": "payment", "latency_ms": 850, "success": false}
		{"app_id": "payment", "latency_ms": 32, "success": true}
	`

	streamReader := strings.NewReader(streamPayload)
	streamDec := json.NewDecoder(streamReader)

	count := 0
	// Loop continuously until we hit EOF
	for {
		var e MetricEvent
		err := streamDec.Decode(&e)

		// Normal exit when we reach the end of the stream
		if err == io.EOF {
			break
		}
		// Hard exit for real parsing errors
		if err != nil {
			log.Fatal("Stream parse error:", err)
		}

		count++
		status := "✅"
		if !e.Success {
			status = "❌"
		}
		fmt.Printf("   Event %d: %s (latency: %dms) %s\n", count, e.AppID, e.Latency, status)
	}

	fmt.Printf("\n   Finished parsing %d objects from stream\n", count)

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - json.NewDecoder(r io.Reader) streams data without loading it all to memory")
	fmt.Println("  - Use it for HTTP Request bodies (req.Body) and very large JSON files")
	fmt.Println("  - Perfect for processing streams of JSON objects (JSONLines format)")
	fmt.Println("  - Loop dec.Decode() until it returns io.EOF")
}
