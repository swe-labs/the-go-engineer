package main

// ============================================================================
// Section 15: Time & Scheduling — Timers & Tickers
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - `time.Timer` for one-off scheduled operations
//   - `time.Ticker` for repeating background intervals
//   - Listening to Timer channels `<-timer.C`
//   - Resource cleanup using `ticker.Stop()`
//
// ENGINEERING DEPTH:
//   A `Timer` or `Ticker` is essentially an OS-level thread sleep mapped to a Go
//   Channel (`<-C`). The Go runtime maintains a global min-heap timer queue.
//   When extreme accuracy is needed, the Go runtime automatically pushes the
//   current `time.Time` into the channel `C` at the targeted hardware tick, unblocking
//   your goroutine. ALWAYS `defer ticker.Stop()`, otherwise the Go Scheduler
//   keeps it in the global min-heap forever, causing massive memory leaks in
//   long-running daemons!
//
// RUN: go run ./15-time-and-scheduling/3-timer-and-ticker
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
			return
		}
	}
}

func timerExample() {

	timer := time.NewTimer(5 * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-timer.C
		fmt.Println("After 1 second")
	}()

	fmt.Println("This is happening inside the main goroutine")

	wg.Wait()

	fmt.Println("program ends")
}
