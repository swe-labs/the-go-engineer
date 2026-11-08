package main

import (
	"runtime"
	"runtime/debug"
	"testing"
)

func TestComputeBudgetReservesHeadroom(t *testing.T) {
	const limit = int64(1) << 30 // 1 GiB
	b, err := ComputeBudget(limit, 100)
	if err != nil {
		t.Fatal(err)
	}
	want := limit * 85 / 100
	if b.MemoryLimit != want {
		t.Fatalf("soft limit = %d, want %d", b.MemoryLimit, want)
	}
	if b.GOMAXPROCS != runtime.GOMAXPROCS(0) {
		t.Fatalf("GOMAXPROCS = %d, want %d", b.GOMAXPROCS, runtime.GOMAXPROCS(0))
	}
}

func TestComputeBudgetSetsRuntimeMemoryLimit(t *testing.T) {
	b, err := ComputeBudget(1<<30, 100)
	if err != nil {
		t.Fatal(err)
	}
	// SetMemoryLimit returns the previous limit and applies the new one. Pass -1
	// to read the current limit, then restore it immediately.
	prev := debug.SetMemoryLimit(-1)
	debug.SetMemoryLimit(prev)
	if prev != b.MemoryLimit {
		t.Fatalf("runtime memory limit = %d, want %d", prev, b.MemoryLimit)
	}
}

func TestComputeBudgetRejectsNonPositiveLimit(t *testing.T) {
	if _, err := ComputeBudget(0, 100); err == nil {
		t.Fatal("expected error for zero container limit")
	}
}

func TestComputeBudgetRejectsNegativeGC(t *testing.T) {
	if _, err := ComputeBudget(1<<30, -1); err == nil {
		t.Fatal("expected error for negative GOGC")
	}
}

func TestParseGCPercentDefault(t *testing.T) {
	got, err := parseGCPercent("")
	if err != nil || got != 100 {
		t.Fatalf("parseGCPercent(\"\") = %d, %v; want 100, nil", got, err)
	}
}

func TestParseGCPercentInvalid(t *testing.T) {
	if _, err := parseGCPercent("fast"); err == nil {
		t.Fatal("expected error for non-numeric GOGC")
	}
}
