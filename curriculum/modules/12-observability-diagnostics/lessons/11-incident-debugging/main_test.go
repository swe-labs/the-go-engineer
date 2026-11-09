package main

import (
	"runtime"
	"sync"
	"testing"
)

func TestGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch := make(chan struct{})
			close(ch)
			<-ch
		}()
	}
	wg.Wait()

	after := runtime.NumGoroutine()
	leaked := after - before
	if leaked > 1 {
		t.Errorf("expected no goroutine leak (diff <= 1), got %d leaked", leaked)
	}
}

func TestActualLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	for i := 0; i < 5; i++ {
		go func() {
			ch := make(chan struct{})
			<-ch
		}()
	}

	runtime.Gosched()
	after := runtime.NumGoroutine()
	leaked := after - before
	if leaked < 5 {
		t.Errorf("expected at least 5 leaked goroutines, got %d", leaked)
	}
}

func BenchmarkMutexContention(b *testing.B) {
	var counter int
	var mu sync.Mutex

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				counter++
				mu.Unlock()
			}()
		}
		wg.Wait()
	}
}

func TestGoroutineProfile(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := make(chan struct{})
		<-ch
	}()

	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	stack := string(buf[:n])

	runtime.Gosched()

	if len(stack) == 0 {
		t.Error("expected non-empty goroutine stack trace")
	}
	if !containsStr(stack, "chan receive") {
		t.Log("note: goroutine may not yet be in chan receive; test is informational")
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
