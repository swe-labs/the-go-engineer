// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Structs
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining named types using the `struct` keyword.
//   - Field access, zero-value initialization, and literal syntax.
//   - Memory layout: field alignment and padding on the stack.
//
// WHY THIS MATTERS:
//   - Structs are the primary mechanism for data grouping in Go.
//     Understanding their memory layout is essential for writing
//     memory-efficient code and correctly interfacing with low-level
//     systems or binary protocols.
//
// RUN:
//   go run ./04-types-design/1-struct
//
// KEY TAKEAWAY:
//   - Structs represent contiguous memory blocks for heterogeneous data.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"fmt"
	"time"
)

// Section 04: Types & Design - Structs
//   - What a struct is and why it exists (grouping related data)
//   - How to define, create, and access struct fields
//   - Zero values for structs (every field gets its type's zero value)
//   - Constructor functions (NewX pattern) for validated creation
//   - Struct pointers: when and why to use them
//   - Struct comparison and copying rules
//

// Server represents a system server with associated metadata.
// Server (Struct): represents a system server with associated metadata.
type Server struct {
	ID       int       // Unique identifier for this server
	Hostname string    // DNS hostname (e.g., "api-prod-01.internal")
	IP       string    // IPv4 address (e.g., "10.0.1.42")
	Region   string    // Cloud region (e.g., "us-east-1")
	CPUCores int       // Number of CPU cores allocated
	MemoryGB int       // RAM in gigabytes
	IsOnline bool      // Current status: true = running, false = stopped
	BootedAt time.Time // When the server was last started
}

// NewServer initializes a Server with validated defaults and current boot time.
// NewServer (Function): initializes a Server with validated defaults and current boot time.
func NewServer(id int, hostname, ip, region string, cpuCores, memoryGB int) Server {
	return Server{
		ID:       id,
		Hostname: hostname,
		IP:       ip,
		Region:   region,
		CPUCores: cpuCores,
		MemoryGB: memoryGB,
		IsOnline: true,
		BootedAt: time.Now(),
	}
}

func main() {
	// 1. Instantiate a struct using literal syntax.
	// Field names are optional but recommended for clarity and forward compatibility.
	webServer := Server{
		ID:       1,
		Hostname: "web-prod-01",
		IP:       "10.0.1.10",
		Region:   "us-east-1",
		CPUCores: 4,
		MemoryGB: 16,
		IsOnline: true,
		BootedAt: time.Now(),
	}

	fmt.Println("=== Structs: Grouping Related Data ===")
	fmt.Println()

	fmt.Printf("Server: %s (%s)\n", webServer.Hostname, webServer.IP)
	fmt.Printf("Region: %s, CPUs: %d, RAM: %dGB\n",
		webServer.Region, webServer.CPUCores, webServer.MemoryGB)
	fmt.Printf("Online: %t, Booted: %s\n",
		webServer.IsOnline, webServer.BootedAt.Format("15:04:05"))
	fmt.Println()

	// 2. Zero-value initialization.
	// Declaring a variable without a literal sets all fields to their zero values.
	var emptyServer Server
	fmt.Printf("Zero value server: %+v\n", emptyServer)
	fmt.Println()

	// 3. Constructor pattern.
	// Using NewServer ensures the object is created in a consistent, valid state.
	dbServer := NewServer(2, "db-prod-01", "10.0.2.20", "us-west-2", 8, 64)
	fmt.Printf("New server: %s (%d cores, %dGB RAM)\n",
		dbServer.Hostname, dbServer.CPUCores, dbServer.MemoryGB)
	fmt.Println()

	dbServer.IsOnline = false
	dbServer.MemoryGB = 128
	fmt.Printf("After upgrade: %s - Online: %t, RAM: %dGB\n",
		dbServer.Hostname, dbServer.IsOnline, dbServer.MemoryGB)
	fmt.Println()

	serverPtr := &dbServer
	serverPtr.IsOnline = true
	fmt.Printf("Via pointer: %s - Online: %t (original modified!)\n",
		dbServer.Hostname, dbServer.IsOnline)
	fmt.Println()

	serverCopy := webServer
	serverCopy.Hostname = "web-staging-01"
	fmt.Printf("Original: %s (unchanged)\n", webServer.Hostname)
	fmt.Printf("Copy:     %s (independent)\n", serverCopy.Hostname)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.2 -> 04-types-design/2-methods")
	fmt.Println("Run    : go run ./04-types-design/2-methods")
	fmt.Println("Current: TI.1 (structs)")
	fmt.Println("---------------------------------------------------")
}
