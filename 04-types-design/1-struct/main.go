// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Structs
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn what a struct is and how to use it to group related data together into a single type.
//
// WHY THIS MATTERS:
//   - Think of a struct like a passport. A passport groups related data about one person: name, nationality, date of birth, photo, passport number. You w...
//
// RUN:
//   go run ./04-types-design/1-struct
//
// KEY TAKEAWAY:
//   - Learn what a struct is and how to use it to group related data together into a single type.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"fmt"
	"time"
)

//
//   - What a struct is and why it exists (grouping related data)
//   - How to define, create, and access struct fields
//   - Zero values for structs (every field gets its type's zero value)
//   - Constructor functions (NewX pattern) for validated creation
//   - Struct pointers: when and why to use them
//   - Struct comparison and copying rules
//

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
	fmt.Println("=== Structs: Grouping Related Data ===")
	fmt.Println()

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

	fmt.Printf("Server: %s (%s)\n", webServer.Hostname, webServer.IP)
	fmt.Printf("Region: %s, CPUs: %d, RAM: %dGB\n",
		webServer.Region, webServer.CPUCores, webServer.MemoryGB)
	fmt.Printf("Online: %t, Booted: %s\n",
		webServer.IsOnline, webServer.BootedAt.Format("15:04:05"))
	fmt.Println()

	var emptyServer Server
	fmt.Printf("Zero value server: %+v\n", emptyServer)
	fmt.Println()

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

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.2 -> 04-types-design/2-methods")
	fmt.Println("Current: TI.1 (structs)")
	fmt.Println("Previous: FE.10 (panic-and-recover)")
	fmt.Println("---------------------------------------------------")
}
