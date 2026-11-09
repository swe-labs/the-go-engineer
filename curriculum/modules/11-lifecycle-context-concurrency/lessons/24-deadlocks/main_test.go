package main

import (
	"sync"
	"testing"
	"time"
)

func TestSafeOrderedLock(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{"same order", func() {
			var mu1, mu2 sync.Mutex
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				mu1.Lock()
				time.Sleep(5 * time.Millisecond)
				mu2.Lock()
				mu2.Unlock()
				mu1.Unlock()
			}()
			go func() {
				defer wg.Done()
				mu1.Lock()
				mu2.Lock()
				mu2.Unlock()
				mu1.Unlock()
			}()
			wg.Wait()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan struct{})
			go func() {
				tt.fn()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(1 * time.Second):
				t.Fatal("deadlock detected: timed out")
			}
		})
	}
}
