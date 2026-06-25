package main

import (
	"sync"
	"testing"
)

func TestWaitGroupsCompiles(t *testing.T) {
	// Validation by compilation and README-driven practice.
}

func TestProcessItems(t *testing.T) {
	items := []WorkItem{
		{ID: 1, Data: "x"},
		{ID: 2, Data: "y"},
	}
	results := processItems(items)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Item.Data != "x" {
		t.Errorf("expected x, got %s", results[0].Item.Data)
	}
}

func TestWaitGroupAddBeforeDone(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{})
	go func() {
		defer wg.Done()
		close(done)
	}()
	<-done
	wg.Wait()
}
