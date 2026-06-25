package main

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func startGoroutine(t *testing.T, fn func()) func() int {
	before := runtime.NumGoroutine()
	fn()
	time.Sleep(10 * time.Millisecond)
	after := runtime.NumGoroutine()
	return func() int { return after - before }
}

func TestNoLeakWithContext(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		wantDiff int
	}{
		{"short timeout", 20 * time.Millisecond, 0},
		{"medium timeout", 50 * time.Millisecond, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := runtime.NumGoroutine()
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			go safeWorker(ctx)
			time.Sleep(100 * time.Millisecond)
			after := runtime.NumGoroutine()
			if leak := after - before; leak > 0 {
				t.Errorf("goroutine leak: %d goroutines remain", leak)
			}
		})
	}
}

func TestLeakyWorker(t *testing.T) {
	done := make(chan bool)
	go leakyWorker(done)
	time.Sleep(10 * time.Millisecond)
	close(done)
}
