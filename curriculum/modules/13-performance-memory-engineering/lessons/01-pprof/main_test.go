package main

import (
	"os"
	"runtime/pprof"
	"testing"
)

func TestCPUProfileCanStartAndStop(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "cpu.pprof")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		t.Fatalf("StartCPUProfile failed: %v", err)
	}
	pprof.StopCPUProfile()

	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("CPU profile file is empty")
	}
}

func TestHeapProfileWritesData(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "heap.pprof")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_ = make([]byte, 1<<20)

	if err := pprof.WriteHeapProfile(f); err != nil {
		t.Fatalf("WriteHeapProfile failed: %v", err)
	}

	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("Heap profile file is empty")
	}
}

func TestProfileGoroutine(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "goroutine.pprof")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	done := make(chan bool)
	go func() {
		<-done
	}()
	defer close(done)

	if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
		t.Fatalf("goroutine profile WriteTo failed: %v", err)
	}

	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("Goroutine profile file is empty")
	}
}

func TestProfileLookupsAreConsistent(t *testing.T) {
	names := []string{"goroutine", "heap", "threadcreate"}
	for _, name := range names {
		p := pprof.Lookup(name)
		if p == nil {
			t.Fatalf("pprof.Lookup(%q) returned nil", name)
		}
	}
}

func TestFibonnaciCorrectness(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{5, 5},
		{10, 55},
		{20, 6765},
	}
	for _, tc := range tests {
		got := fibonacci(tc.input)
		if got != tc.expected {
			t.Errorf("fibonacci(%d) = %d; want %d", tc.input, got, tc.expected)
		}
	}
}
