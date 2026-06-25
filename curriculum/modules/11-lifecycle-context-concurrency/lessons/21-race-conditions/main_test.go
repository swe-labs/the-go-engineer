package main

import (
	"sync"
	"testing"
)

func TestSafeCounter(t *testing.T) {
	tests := []struct {
		name            string
		goroutines      int
		opsPerGoroutine int
	}{
		{"single goroutine", 1, 10},
		{"two goroutines", 2, 100},
		{"ten goroutines", 10, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SafeCounter{}
			var wg sync.WaitGroup
			total := tt.goroutines * tt.opsPerGoroutine
			wg.Add(total)
			for i := 0; i < total; i++ {
				go func() {
					defer wg.Done()
					c.Add(1)
				}()
			}
			wg.Wait()
			if got := c.Value(); got != total {
				t.Errorf("SafeCounter.Value() = %d, want %d", got, total)
			}
		})
	}
}

func TestRacyCounterValue(t *testing.T) {
	c := &RacyCounter{}
	c.Add(5)
	if got := c.Value(); got != 5 {
		t.Errorf("RacyCounter.Value() = %d, want 5", got)
	}
}
