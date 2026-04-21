// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Stage 07: Concurrency — Select Deep Dive
// Level: Advanced
// ============================================================================
//
// RUN: go run ./07-concurrency/01-concurrency/goroutines/9-select-deep-dive
// ============================================================================

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Select Statement Deep Dive ===")
	fmt.Println()

	basicSelect()
	timeoutPattern()
	nonBlockingSelect()
	contextCancellation()
	fanInPattern()

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: GC.10 sync primitives")
	fmt.Println("   Current: GC.9 (select deep dive)")
	fmt.Println("---------------------------------------------------")
}

func basicSelect() {
	fmt.Println("--- 1. Basic Select ---")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("  Received: %s\n", msg)
		case msg := <-ch2:
			fmt.Printf("  Received: %s\n", msg)
		}
	}
	fmt.Println()
}

func timeoutPattern() {
	fmt.Println("--- 2. Timeout Pattern ---")

	slowOperation := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		slowOperation <- "done"
	}()

	select {
	case result := <-slowOperation:
		fmt.Printf("  Got result: %s\n", result)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("  ⏰ Operation timed out after 500ms")
	}
	fmt.Println()
}

func nonBlockingSelect() {
	fmt.Println("--- 3. Non-blocking Select ---")

	ch := make(chan int, 1)

	select {
	case val := <-ch:
		fmt.Printf("  Received: %d\n", val)
	default:
		fmt.Println("  Channel empty — no blocking!")
	}

	ch <- 42
	select {
	case ch <- 100:
		fmt.Println("  Sent 100")
	default:
		fmt.Println("  Channel full — no blocking!")
	}
	fmt.Println()
}

func contextCancellation() {
	fmt.Println("--- 4. Context Cancellation ---")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	resultCh := make(chan string)
	go func() {
		time.Sleep(500 * time.Millisecond)
		resultCh <- "work complete"
	}()

	select {
	case result := <-resultCh:
		fmt.Printf("  Result: %s\n", result)
	case <-ctx.Done():
		fmt.Printf("  ❌ Cancelled: %v\n", ctx.Err())
	}
	fmt.Println()
}

func fanInPattern() {
	fmt.Println("--- 5. Fan-In Pattern ---")

	producers := make([]<-chan string, 3)
	for i := 0; i < 3; i++ {
		producers[i] = produce(i)
	}

	merged := fanIn(producers...)

	for i := 0; i < 6; i++ { // 3 producers × 2 messages each
		fmt.Printf("  %s\n", <-merged)
	}
	fmt.Println()
}

func produce(id int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := 0; i < 2; i++ {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			ch <- fmt.Sprintf("Producer %d: message %d", id, i)
		}
	}()
	return ch
}

func fanIn(channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for msg := range c {
				out <- msg
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
