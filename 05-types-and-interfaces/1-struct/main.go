// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"time"
)

// ============================================================================
// Section 5: Types & Interfaces — Structs
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What a struct is and why it exists (grouping related data)
//   - How to define, create, and access struct fields
//   - Zero values for structs (every field gets its type's zero value)
//   - Constructor functions (NewX pattern) for validated creation
//   - Struct pointers: when and why to use them
//   - Struct comparison and copying rules
//
// ANALOGY:
//   Think of a struct like a passport. A passport groups related data
//   about one person: name, nationality, date of birth, photo, passport number.
//   You wouldn't scatter this data across 6 separate variables — you'd put
//   it in one structured document. That's exactly what a struct does in code.
//
// RUN: go run ./05-types-and-interfaces/1-struct
// ============================================================================

// Server represents a cloud server instance.
// Each field has a name, a type, and an optional tag (covered in Section 11: Encoding).
//
// MEMORY LAYOUT: Go lays out struct fields contiguously in memory.
// The compiler may add "padding" bytes between fields for alignment.
// Ordering fields from largest to smallest can reduce memory usage.
// (This is an advanced optimization — don't worry about it as a beginner.)
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

// NewServer is a CONSTRUCTOR FUNCTION — Go's alternative to class constructors.
//
// Go doesn't have classes or "new" keywords. Instead, by convention:
//   - Functions that create types are named New<Type> (e.g., NewServer, NewReader)
//   - They validate inputs and set sensible defaults
//   - They return the struct (or a pointer to it)
//
// WHY USE CONSTRUCTORS?
//
//	Without a constructor, someone could write:
//	  s := Server{}  ← All zero values. Hostname is "", CPUCores is 0.
//	That's a broken server. A constructor ensures valid state from creation.
func NewServer(id int, hostname, ip, region string, cpuCores, memoryGB int) Server {
	return Server{
		ID:       id,
		Hostname: hostname,
		IP:       ip,
		Region:   region,
		CPUCores: cpuCores,
		MemoryGB: memoryGB,
		IsOnline: true,       // New servers start online
		BootedAt: time.Now(), // Record boot time automatically
	}
}

func main() {
	fmt.Println("=== Structs: Grouping Related Data ===")
	fmt.Println()

	// --- METHOD 1: Named field initialization (RECOMMENDED) ---
	// Every field is explicitly named. Order doesn't matter.
	// Any field you omit gets its zero value.
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

	// Access fields with dot notation: structVar.FieldName
	fmt.Printf("Server: %s (%s)\n", webServer.Hostname, webServer.IP)
	fmt.Printf("Region: %s, CPUs: %d, RAM: %dGB\n",
		webServer.Region, webServer.CPUCores, webServer.MemoryGB)
	fmt.Printf("Online: %t, Booted: %s\n",
		webServer.IsOnline, webServer.BootedAt.Format("15:04:05"))
	fmt.Println()

	// --- METHOD 2: Zero value struct ---
	// Declaring without initializing gives you the zero value for EVERY field:
	//   int → 0, string → "", bool → false, time.Time → 0001-01-01
	var emptyServer Server
	fmt.Printf("Zero value server: %+v\n", emptyServer)
	// %+v prints field names AND values — very useful for debugging structs.
	fmt.Println()

	// --- METHOD 3: Constructor function ---
	// The constructor sets defaults (IsOnline=true, BootedAt=now).
	// You can't forget to set critical fields because they're parameters.
	dbServer := NewServer(2, "db-prod-01", "10.0.2.20", "us-west-2", 8, 64)
	fmt.Printf("New server: %s (%d cores, %dGB RAM)\n",
		dbServer.Hostname, dbServer.CPUCores, dbServer.MemoryGB)
	fmt.Println()

	// --- MODIFYING FIELDS ---
	// Struct fields are mutable — you can change them after creation.
	dbServer.IsOnline = false
	dbServer.MemoryGB = 128 // Upgraded RAM
	fmt.Printf("After upgrade: %s — Online: %t, RAM: %dGB\n",
		dbServer.Hostname, dbServer.IsOnline, dbServer.MemoryGB)
	fmt.Println()

	// --- STRUCT POINTERS ---
	// A pointer to a struct lets you modify the ORIGINAL struct.
	// Without a pointer, you get a COPY (same as arrays in Section 03).
	//
	// Go automatically dereferences struct pointers:
	//   serverPtr.Hostname is the same as (*serverPtr).Hostname
	//   You almost never need to write (*ptr).Field explicitly.
	serverPtr := &dbServer    // & creates a pointer to dbServer
	serverPtr.IsOnline = true // Modifies the ORIGINAL dbServer
	fmt.Printf("Via pointer: %s — Online: %t (original modified!)\n",
		dbServer.Hostname, dbServer.IsOnline)
	fmt.Println()

	// --- STRUCT COPYING ---
	// Assigning one struct to another creates a FULL COPY.
	// Modifying the copy does NOT affect the original.
	serverCopy := webServer
	serverCopy.Hostname = "web-staging-01"
	fmt.Printf("Original: %s (unchanged)\n", webServer.Hostname)
	fmt.Printf("Copy:     %s (independent)\n", serverCopy.Hostname)

	// KEY TAKEAWAY:
	// - Structs group related data into a single type (like a passport or record)
	// - Use named fields for clarity: Server{Hostname: "x"} not Server{"x"}
	// - Use NewX() constructors to ensure valid initial state
	// - Structs are value types: assignment creates a COPY
	// - Use pointers (&server) when you need to modify the original
}
