package main

import (
	"context"
	"fmt"
	"time"
)

func leakyWorker(done chan bool) {
	// Blocks forever waiting on a channel that never sends
	<-done
}

func safeWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	done := make(chan bool)
	go leakyWorker(done)
	fmt.Println("Leaky worker started (will leak)")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	go safeWorker(ctx)

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Safe worker exited via context cancellation")
	fmt.Println("\nTip: use `runtime.NumGoroutine()` to detect leaks in tests")
}
