package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 07: Concurrency - Why concurrency exists
//
// Run: go run ./07-concurrency/01-concurrency/goroutines/0-why-concurrency-exists

func main() {
	fmt.Println("=== GC.0 Why concurrency exists ===")
	fmt.Println("Understand why overlapping waits is the reason concurrency exists in everyday backend code.")
	fmt.Println()
	fmt.Println("- Blocking I/O creates idle time on one execution path.")
	fmt.Println("- Goroutines let another task make progress during that wait.")
	fmt.Println("- Concurrency improves throughput only when there is real waiting to overlap.")
	fmt.Println()
	fmt.Println("Most network services are dominated by waits. Concurrency helps most when a request fan-out, queue, socket, or disk call spends time blocked.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GC.1")
	fmt.Println("Current: GC.0 (why concurrency exists)")
	fmt.Println("---------------------------------------------------")
}
