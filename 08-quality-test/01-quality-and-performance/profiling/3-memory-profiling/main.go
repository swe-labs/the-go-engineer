package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 08: Quality & Testing - Memory profiling
//
// Run: go run ./08-quality-test/01-quality-and-performance/profiling/3-memory-profiling

func main() {
	fmt.Println("=== PR.3 Memory profiling ===")
	fmt.Println("Learn how heap profiles reveal where memory is retained and where allocation pressure is coming from.")
	fmt.Println()
	fmt.Println("- Profiles show where memory is retained.")
	fmt.Println("- Allocation rate and live heap are related but different signals.")
	fmt.Println("- Measure before changing data structures or allocation patterns.")
	fmt.Println()
	fmt.Println("Memory tuning starts with visibility. Without a profile, you are usually guessing at the wrong thing.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: PR.4")
	fmt.Println("Current: PR.3 (memory profiling)")
	fmt.Println("---------------------------------------------------")
}
