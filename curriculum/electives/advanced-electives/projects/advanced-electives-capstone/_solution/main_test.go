package main

import (
	"sync"
	"testing"
)

func TestRingBufferConcurrency(t *testing.T) {
	rb := NewRingBuffer[int](100)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			rb.Push(val)
		}(i)
	}
	wg.Wait()

	if rb.Len() != 50 {
		t.Fatalf("expected 50 items, got %d", rb.Len())
	}
}
