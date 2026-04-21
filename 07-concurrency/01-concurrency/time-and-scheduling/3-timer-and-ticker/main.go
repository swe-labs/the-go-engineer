// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Stage 07: Time & Scheduling - Timers & Tickers
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - `time.Timer` for one-off scheduled operations
//   - `time.Ticker` for repeating background intervals
//   - Listening to timer channels with `<-timer.C`
//   - Resource cleanup using `ticker.Stop()`
//
// ENGINEERING DEPTH:
//   A `Timer` or `Ticker` is backed by the Go runtime timer heap.
//   When the target time arrives, the runtime delivers a value on the timer's
//   channel and unblocks your goroutine. Always `defer ticker.Stop()`, or the
//   runtime will keep the timer alive longer than necessary.
//
// RUN: go run ./07-concurrency/01-concurrency/time-and-scheduling/3-timer-and-ticker
// ============================================================================

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	timerExample()
	fmt.Println()

	ticker := time.NewTicker(1 * time.Second)
	counter := 0
	defer ticker.Stop()

	for range ticker.C {
		counter++
		fmt.Println("Tick")
		if counter >= 5 {
			fmt.Println("stopped")
			break
		}
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TM.7 console reminder")
	fmt.Println("   Current: TM.3 (timers & tickers)")
	fmt.Println("---------------------------------------------------")
}

func timerExample() {
	timer := time.NewTimer(5 * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-timer.C
		fmt.Println("After 5 seconds")
	}()

	fmt.Println("This is happening inside the main goroutine")

	wg.Wait()

	fmt.Println("program ends")
}
