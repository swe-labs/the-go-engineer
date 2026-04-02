// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"time"
)

// ============================================================================
// Section 9: Concurrency — Channels
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What channels are: typed communication pipes between goroutines
//   - The channel axiom: "Don't communicate by sharing memory.
//     Share memory by communicating."
//   - Sending and receiving: ch <- value (send), value := <-ch (receive)
//   - Blocking behavior: unbuffered channels synchronize sender and receiver
//   - Channel direction: send-only (chan<-) and receive-only (<-chan)
//   - Practical patterns: result collection, fan-out
//
// ANALOGY:
//   A channel is like a mailbox between two neighbors.
//   - Neighbor A (sender) puts a letter in the mailbox: ch <- letter
//   - Neighbor B (receiver) takes the letter out: letter := <-ch
//
//   With an UNBUFFERED channel (no mailbox space), A must WAIT at the
//   mailbox until B arrives to take the letter. They synchronize.
//
//   With a BUFFERED channel (mailbox with slots), A can drop letters
//   and leave — until the mailbox is full. Then A waits too.
//
// ENGINEERING DEPTH:
//   Channels are implemented as a struct (hchan) containing:
//     - A circular buffer (for buffered channels)
//     - A mutex for thread-safe access
//     - Two wait queues: one for blocked senders, one for blocked receivers
//   When a goroutine blocks on a channel, the Go scheduler parks it
//   (removes it from the OS thread) and schedules another goroutine.
//   This is a userspace context switch — ~10ns vs ~1μs for OS threads.
//
// RUN: go run ./09-concurrency/3-channels
// ============================================================================

// ScanResult represents the result of scanning a server port.
// This is the kind of data you'd send between goroutines in a real app.
type ScanResult struct {
	Host   string // Target hostname or IP
	Port   int    // Port number scanned
	IsOpen bool   // Whether the port responded
}

// scanPort simulates scanning a network port on a host.
// It sends the result through a channel instead of returning it.
//
// CHANNEL DIRECTION: The parameter type is chan<- ScanResult (SEND-ONLY).
// This means the function can only SEND into the channel, never receive.
// This is a compile-time safety feature:
//
//	chan<- T = send-only channel
//	<-chan T = receive-only channel
//	chan T   = bidirectional channel (can send and receive)
func scanPort(host string, port int, results chan<- ScanResult) {
	// Simulate network latency (different ports take different time)
	delay := time.Duration(port%3+1) * 100 * time.Millisecond
	time.Sleep(delay)

	// Simulate: ports 80 and 443 are open, others are closed
	isOpen := port == 80 || port == 443 || port == 22

	// Send the result into the channel.
	// This BLOCKS until someone reads from the other end.
	// That's the synchronization magic of unbuffered channels.
	results <- ScanResult{
		Host:   host,
		Port:   port,
		IsOpen: isOpen,
	}
}

func main() {
	fmt.Println("=== Channels: Communication Between Goroutines ===")
	fmt.Println()

	// =====================================================================
	// EXAMPLE 1: Basic send and receive
	// =====================================================================
	fmt.Println("--- Basic Channel ---")

	// make(chan Type) creates an UNBUFFERED channel.
	// It can transport values of the specified type between goroutines.
	greetings := make(chan string)

	// Send a value in a goroutine.
	// We MUST send in a goroutine because unbuffered channels BLOCK the sender
	// until a receiver is ready. If we sent in main, main would block forever
	// (deadlock) because no one is receiving yet.
	go func() {
		greetings <- "Hello from a goroutine!" // Send (blocks until received)
	}()

	// Receive the value in main.
	// <-greetings blocks until a sender puts a value in the channel.
	msg := <-greetings // Receive (blocks until sent)
	fmt.Printf("  Received: %s\n", msg)
	fmt.Println()

	// =====================================================================
	// EXAMPLE 2: Collecting results from multiple goroutines
	// =====================================================================
	fmt.Println("--- Port Scanner (Concurrent Result Collection) ---")

	// Create a channel to collect scan results.
	results := make(chan ScanResult)

	// Define the ports to scan
	portsToScan := []int{22, 80, 443, 3306, 5432, 8080}
	host := "10.0.1.42"

	// Launch a goroutine for EACH port scan.
	// All scans run CONCURRENTLY — much faster than scanning sequentially.
	for _, port := range portsToScan {
		go scanPort(host, port, results) // Each goroutine sends its result to the channel
	}

	// Collect ALL results from the channel.
	// We know exactly how many results to expect (one per port),
	// so we read exactly len(portsToScan) values.
	fmt.Printf("  Scanning %s on %d ports...\n\n", host, len(portsToScan))
	for i := 0; i < len(portsToScan); i++ {
		result := <-results // Receive one result (blocks until available)

		status := "❌ closed"
		if result.IsOpen {
			status = "✅ OPEN"
		}
		fmt.Printf("  %s:%d → %s\n", result.Host, result.Port, status)
	}

	fmt.Println()

	// =====================================================================
	// EXAMPLE 3: Channels as synchronization (signaling done)
	// =====================================================================
	fmt.Println("--- Channel as Signal (done pattern) ---")

	// A channel can carry no data — just signal that something happened.
	// struct{} takes ZERO bytes of memory — the most efficient signal type.
	done := make(chan struct{})

	go func() {
		fmt.Println("  Background task: processing...")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("  Background task: complete!")
		close(done) // Closing a channel signals ALL receivers
	}()

	<-done // Block until the channel is closed
	fmt.Println("  Main: received done signal, continuing")

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Channels are typed pipes: make(chan string), make(chan int)")
	fmt.Println("  - Send: ch <- value | Receive: value := <-ch")
	fmt.Println("  - Unbuffered channels SYNCHRONIZE sender and receiver")
	fmt.Println("  - Use chan<- (send-only) and <-chan (receive-only) for safety")
	fmt.Println("  - close(ch) signals all receivers that no more values are coming")
	fmt.Println("  - Channels replace shared memory + mutexes for most use cases")
}
