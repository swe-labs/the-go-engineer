// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"time"
)

// ============================================================================
// Stage 07: Concurrency - Channels
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What channels are: typed communication pipes between goroutines
//   - The channel axiom: "Do not communicate by sharing memory. Share memory by communicating."
//   - Sending and receiving: ch <- value (send), value := <-ch (receive)
//   - Blocking behavior: unbuffered channels synchronize sender and receiver
//   - Channel direction: send-only (chan<-) and receive-only (<-chan)
//   - Practical patterns: result collection and fan-out
//
// ANALOGY:
//   A channel is like a mailbox between two neighbors.
//   - Neighbor A (sender) puts a letter in the mailbox: ch <- letter
//   - Neighbor B (receiver) takes the letter out: letter := <-ch
//
//   With an unbuffered channel (no mailbox space), A must wait at the mailbox
//   until B arrives to take the letter. They synchronize.
//
//   With a buffered channel (mailbox with slots), A can drop letters there
//   until the mailbox is full. Then A waits too.
//
// ENGINEERING DEPTH:
//   Channels are implemented as a struct (hchan) containing:
//     - A circular buffer (for buffered channels)
//     - A mutex for thread-safe access
//     - Two wait queues: one for blocked senders, one for blocked receivers
//   When a goroutine blocks on a channel, the Go scheduler parks it and lets
//   another goroutine run. This is a user-space context switch instead of a
//   full OS-thread handoff.
//
// RUN: go run ./07-concurrency/01-concurrency/goroutines/3-channels
// ============================================================================

type ScanResult struct {
	Host   string
	Port   int
	IsOpen bool
}

func scanPort(host string, port int, results chan<- ScanResult) {
	delay := time.Duration(port%3+1) * 100 * time.Millisecond
	time.Sleep(delay)

	isOpen := port == 80 || port == 443 || port == 22

	results <- ScanResult{
		Host:   host,
		Port:   port,
		IsOpen: isOpen,
	}
}

func main() {
	fmt.Println("=== Channels: Communication Between Goroutines ===")
	fmt.Println()

	fmt.Println("--- Basic Channel ---")
	greetings := make(chan string)
	go func() {
		greetings <- "Hello from a goroutine!"
	}()

	msg := <-greetings
	fmt.Printf("  Received: %s\n", msg)
	fmt.Println()

	fmt.Println("--- Port Scanner (Concurrent Result Collection) ---")
	results := make(chan ScanResult)
	portsToScan := []int{22, 80, 443, 3306, 5432, 8080}
	host := "10.0.1.42"

	for _, port := range portsToScan {
		go scanPort(host, port, results)
	}

	fmt.Printf("  Scanning %s on %d ports...\n\n", host, len(portsToScan))
	for i := 0; i < len(portsToScan); i++ {
		result := <-results

		status := "[closed]"
		if result.IsOpen {
			status = "[OPEN]"
		}
		fmt.Printf("  %s:%d -> %s\n", result.Host, result.Port, status)
	}

	fmt.Println()
	fmt.Println("--- Channel as Signal (done pattern) ---")
	done := make(chan struct{})

	go func() {
		fmt.Println("  Background task: processing...")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("  Background task: complete!")
		close(done)
	}()

	<-done
	fmt.Println("  Main: received done signal, continuing")

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Channels are typed pipes: make(chan string), make(chan int)")
	fmt.Println("  - Send: ch <- value | Receive: value := <-ch")
	fmt.Println("  - Unbuffered channels synchronize sender and receiver")
	fmt.Println("  - Use chan<- (send-only) and <-chan (receive-only) for safety")
	fmt.Println("  - close(ch) signals all receivers that no more values are coming")
	fmt.Println("  - Channels replace shared memory + mutexes for many coordination cases")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: GC.4 buffered channels")
	fmt.Println("   Current: GC.3 (channels (unbuffered))")
	fmt.Println("---------------------------------------------------")
}
