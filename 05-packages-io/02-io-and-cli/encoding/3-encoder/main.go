// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ============================================================================
// Stage 05: Encoding — JSON Encoder (Streaming)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - json.NewEncoder: streaming JSON directly to an io.Writer (like a file or HTTP response)
//   - Encoder vs Marshal: when to use which
//   - Memory efficiency for large payloads
//
// ANALOGY:
//   json.Marshal is like writing an entire letter at your desk, putting it in an envelope,
//   and then carrying the envelope to the post office. (High memory, simple).
//
//   json.NewEncoder is like sitting AT the post office and writing the letter directly
//   onto the conveyor belt. (Low memory, efficient for large data).
//
// ENGINEERING DEPTH:
//   When building high-throughput microservices, allocating a massive `[]byte`
//   for every single HTTP response via `json.Marshal` will cause devastating
//   Garbage Collection pausing. By giving `json.NewEncoder` an `http.ResponseWriter`
//   (which implements `io.Writer`), the JSON encoder writes the encoded bytes
//   directly into the open TCP network socket in small, highly-optimized chunks.
//   The memory footprint stays nearly flat regardless of the payload size.
//
// WHEN TO USE WHICH?
//   - Use json.Marshal when you need the []byte (e.g., storing in a database or Redis).
//   - Use json.NewEncoder when writing directly to an io.Writer (e.g., HTTP ResponseWriter, file).
//
// RUN: go run ./05-packages-io/02-io-and-cli/encoding/3-encoder
// ============================================================================

type DeviceLog struct {
	DeviceID  string  `json:"device_id"`
	Timestamp int64   `json:"timestamp"`
	CPUUsage  float64 `json:"cpu_usage"`
	MemUsage  float64 `json:"mem_usage"`
	Status    string  `json:"status"`
}

func main() {
	fmt.Println("=== JSON Encoder (Streaming) ===")
	fmt.Println()

	// --- 1. Streaming directly to Stdout ---
	// Create an encoder that writes directly to standard output (the terminal).
	// We don't need to load the JSON into a byte slice first.
	fmt.Println("1️⃣  Streaming to os.Stdout:")
	enc := json.NewEncoder(os.Stdout)

	// SetIndent works just like MarshalIndent
	enc.SetIndent("   ", "  ")

	sampleLog := DeviceLog{
		DeviceID:  "server-prod-01",
		Timestamp: 1684501234,
		CPUUsage:  82.5,
		MemUsage:  65.0,
		Status:    "WARNING",
	}

	// Encode writes the JSON representation to the stream
	if err := enc.Encode(&sampleLog); err != nil {
		log.Fatal("Encode failed:", err)
	}
	fmt.Println()

	// --- 2. Streaming to a File ---
	// This is where NewEncoder shines. If you have 100,000 logs,
	// you can stream them directly to disk without loading them all in memory.
	fmt.Println("2️⃣  Streaming to a file (memory-efficient):")

	file, err := os.CreateTemp("", "logs-*.jsonl")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name()) // Clean up after demo
	defer file.Close()

	// Create an encoder writing to the file
	fileEnc := json.NewEncoder(file)

	// Simulate streaming multiple records (JSONL format - JSON Lines)
	for i := 1; i <= 3; i++ {
		logData := DeviceLog{
			DeviceID:  fmt.Sprintf("sensor-%02d", i),
			Timestamp: 1684501200 + int64(i*10),
			Status:    "OK",
		}

		// Encode writes exactly one JSON object followed by a newline (\n).
		// This is perfect for appending structural logs.
		if err := fileEnc.Encode(logData); err != nil {
			log.Fatal("File encode failed:", err)
		}
	}

	fmt.Printf("   Successfully streamed 3 JSON records to: %s\n", file.Name())

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - json.NewEncoder(w io.Writer) streams data directly to the writer")
	fmt.Println("  - Highly memory-efficient for large files or HTTP responses")
	fmt.Println("  - Use this for APIs (ResponseWriter) instead of json.Marshal")
	fmt.Println("  - Encode() automatically appends a newline (\n) after the JSON")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: EN.4 JSON decoder (stream)")
	fmt.Println("   Current: EN.3 (JSON encoder (stream))")
	fmt.Println("---------------------------------------------------")
}
