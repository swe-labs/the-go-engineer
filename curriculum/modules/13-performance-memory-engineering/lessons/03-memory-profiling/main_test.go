package main

import (
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
)

func TestLeakIncreasesMemory(t *testing.T) {
	var before, after runtime.MemStats

	var holder LeakHolder

	runtime.ReadMemStats(&before)
	for i := 0; i < 20; i++ {
		holder.leak(1 << 20)
	}
	runtime.ReadMemStats(&after)

	leaked := after.HeapAlloc - before.HeapAlloc
	if leaked < 10<<20 {
		t.Errorf("Expected at least 10MB leaked, got %d bytes", leaked)
	}
}

func TestNoLeakDoesNotIncreaseMemory(t *testing.T) {
	var before, after runtime.MemStats

	var holder LeakHolder

	runtime.ReadMemStats(&before)
	for i := 0; i < 20; i++ {
		holder.noLeak(1 << 20)
	}
	runtime.GC()
	runtime.ReadMemStats(&after)

	leaked := int64(after.HeapAlloc) - int64(before.HeapAlloc)
	if leaked > 1<<20 {
		t.Errorf("Expected no leak (under 1MB), got %d bytes", leaked)
	}
}

func TestHeapProfileWritten(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "heap.pprof")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var holder LeakHolder
	holder.leak(5 << 20)

	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		t.Fatal(err)
	}

	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("heap profile is empty")
	}
}

func TestLeakTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		iter     int
		size     int
		wantLeak bool
	}{
		{"single leak", 1, 1 << 20, true},
		{"multiple leaks", 5, 2 << 20, true},
		{"small alloc", 100, 64, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var before, after runtime.MemStats
			var holder LeakHolder

			runtime.ReadMemStats(&before)
			for i := 0; i < tc.iter; i++ {
				holder.leak(tc.size)
			}
			runtime.ReadMemStats(&after)

			diff := after.HeapAlloc - before.HeapAlloc
			if tc.wantLeak && diff < uint64(tc.iter*tc.size/2) {
				t.Errorf("expected leak of ~%d bytes, got %d", tc.iter*tc.size, diff)
			}
		})
	}
}
