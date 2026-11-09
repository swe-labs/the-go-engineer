package main

import (
	"sync"
	"testing"
)

func TestRaceDetectorSafeCounter(t *testing.T) {
	// Run with: go test -race .
	tests := []struct {
		name       string
		goroutines int
	}{
		{"concurrent_10", 10},
		{"concurrent_100", 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SharedCounter{}
			var wg sync.WaitGroup
			wg.Add(tt.goroutines)
			for i := 0; i < tt.goroutines; i++ {
				go func() {
					defer wg.Done()
					c.Add(1)
				}()
			}
			wg.Wait()
			if got := c.Value(); got != tt.goroutines {
				t.Errorf("got %d, want %d", got, tt.goroutines)
			}
		})
	}
}
