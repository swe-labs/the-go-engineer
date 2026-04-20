package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 08: Quality & Testing - Why memory layout matters
//
// Run: go run ./08-quality-test/01-quality-and-performance/profiling/6-memory-layout

func main() {
	fmt.Println("=== PR.6 Why memory layout matters ===")
	fmt.Println("Understand why field order and access pattern influence cache behavior and practical cost.")
	fmt.Println()
	fmt.Println("- Field order changes struct size.")
	fmt.Println("- Cache locality changes traversal cost.")
	fmt.Println("- Compact hot data is easier for the CPU to serve quickly.")
	fmt.Println()
	fmt.Println("Layout tuning is a hot-path tool, not a replacement for good algorithms.")
}
